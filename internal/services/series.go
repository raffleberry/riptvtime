package services

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/raffleberry/riptvtime/internal/db"
	"github.com/raffleberry/riptvtime/internal/meta"
	"github.com/raffleberry/riptvtime/internal/services/state"
	"golang.org/x/sync/errgroup"
)

var (
	ErrInvalidData = errors.New("Invalid Data")
	ErrNotFound    = errors.New("Not Found")
	ErrDuplicate   = errors.New("Duplicate already in database")
)

type SeriesService struct {
	db   db.Db
	meta meta.Meta
	ipt  *ImportSvc
}

func NewTvService(db db.Db, meta meta.Meta, ipt *ImportSvc) *SeriesService {
	return &SeriesService{
		db:   db,
		meta: meta,
		ipt:  ipt,
	}
}

func (srv *SeriesService) Search(searchTerm string, page int) (*SeriesSearchResult, error) {

	metaRes, err := srv.meta.Search(searchTerm, page)
	if err != nil {
		return nil, err
	}

	rv := SeriesSearchResult{
		Page:         metaRes.Page,
		TotalPages:   metaRes.TotalPages,
		TotalResults: metaRes.TotalResults,
	}

	for _, r := range metaRes.Results {

		status, err := srv.db.SeriesStatusGet(r.Id)
		if err != nil {
			return nil, err
		}

		item := SeriesSearchItem{
			TvSearchResult: r,
			Status:         status,
		}

		rv.Results = append(rv.Results, item)

	}

	return &rv, nil
}

func (srv *SeriesService) GetTvCacheExpireTime(res *meta.TvDetails) time.Time {
	rv := db.GetInProdExpireTime()
	if !res.NextEpisodeToAir.AirDate.IsZero() && res.NextEpisodeToAir.AirDate.Before(rv) {
		rv = res.NextEpisodeToAir.AirDate
	}
	return rv
}

// Cached
func (srv *SeriesService) GetTvMeta(mId int) (*meta.TvDetails, error) {
	what := "TvDetails"
	key := strconv.Itoa(mId)
	var rv *meta.TvDetails

	var refresh = func() (*db.Cached, error) {
		var err error
		rv, err = srv.meta.GetTvDetails(mId)
		if err != nil {
			return nil, err
		}
		jsonStr, err := json.Marshal(rv)
		if err != nil {
			return nil, err
		}

		expireTime := time.Now()
		if rv.InProduction {
			expireTime = srv.GetTvCacheExpireTime(rv)
		} else {
			expireTime = db.GetNotInProdExpireTime()
		}

		rv := &db.Cached{
			What:      what,
			Key:       key,
			JsonData:  string(jsonStr),
			ExpiredAt: expireTime,
		}

		err = srv.db.CacheSet(rv)
		if err != nil {
			return nil, err
		}
		return rv, nil
	}

	cd, err := srv.db.CacheGet(what, key)

	if errors.Is(err, db.ErrNotFound) {
		cd, err = refresh()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if cd.ExpiredAt.IsZero() {
		// 1. old transitional logic after migration
		cd.ExpiredAt = cd.UpdatedAt.Add(48 * time.Hour)
	}

	if time.Now().After(cd.ExpiredAt) {
		slog.Debug("Cache expired, refreshing", "expiredAt", cd.ExpiredAt, "UpdatedAt", cd.UpdatedAt)
		cd, err = refresh()
		if err != nil {
			return nil, err
		}
	}

	if rv == nil {
		err = json.Unmarshal([]byte(cd.JsonData), &rv)
		if err != nil {
			return nil, err
		}
	}

	return rv, nil

}

// Get Fresh Tv Show Details (Fresh enough)
func (srv *SeriesService) GetDetails(mId int, withEpsDetails bool) (*SeriesFullItem, error) {

	res, err := srv.GetTvMeta(mId)
	if err != nil {
		return nil, err
	}

	epsAired := 0

	for i, s := range res.Seasons {

		if isLegitSeason(s.Name, s.SeasonNumber, res.NumberOfSeasons) {
			if s.SeasonNumber < res.LastEpisodeToAir.SeasonNumber {
				epsAired += s.EpisodeCount
			} else if s.SeasonNumber == res.LastEpisodeToAir.SeasonNumber {
				epsAired += res.LastEpisodeToAir.EpisodeNumber
			}
		}

		if withEpsDetails {

			for epNo := 1; epNo <= s.EpisodeCount; epNo++ {
				ep, err := srv.getEpisodeDetails(mId, s.SeasonNumber, epNo)
				if err != nil {
					if !isLegitSeason(s.Name, s.SeasonNumber, res.NumberOfSeasons) {
						continue
					}
					return nil, err
				}

				res.Seasons[i].Episodes = append(res.Seasons[i].Episodes, meta.TvEpisode{
					Id:            int(ep.MId),
					Name:          ep.Name,
					Overview:      ep.Overview,
					SeasonNumber:  ep.Season,
					EpisodeNumber: ep.Episode,
					AirDate:       ep.AirDate,
					Runtime:       ep.Runtime,
					MName:         ep.MName,
				})
			}
		}
	}

	epsWatched := []SeriesEpisode{}

	tEps, err := srv.db.SeriesTrackedEps(mId)
	if err != nil {
		return nil, err
	}

	for _, tep := range tEps {
		idx := slices.IndexFunc(epsWatched, func(ep SeriesEpisode) bool {
			return ep.S == tep.Season && ep.E == tep.Episode
		})
		if idx == -1 {
			epsWatched = append(epsWatched, SeriesEpisode{
				S:         tep.Season,
				E:         tep.Episode,
				Cnt:       1,
				CreatedAt: tep.CreatedAt,
			})
		} else {
			epsWatched[idx].Cnt += 1
		}
	}

	return &SeriesFullItem{
		res,
		epsAired,
		epsWatched,
	}, nil
}

func (srv *SeriesService) TrackedAll() (*[]db.TvSeries, error) {
	rv, err := srv.db.SeriesTrackedAll()
	if err != nil {
		return nil, err
	}

	for i := range *rv {
		srs := &(*rv)[i]
		srs.TrackingStatus, err = srv.deriveStatus(int(srs.MId), srs.TrackingStatus)
	}

	return rv, nil
}

func (srv *SeriesService) freshSeriesData(series []db.TvSeries) ([]*SeriesFullItem, error) {
	fsd := []*SeriesFullItem{}

	var mu sync.Mutex

	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(10)

	for _, srs := range series {

		g.Go(func() error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			res, err := srv.GetDetails(int(srs.MId), false)
			if err != nil {
				slog.Error("Error while fetching data", "error", err, "mId", srs.MId)
				return err
			}

			go func() {
				if res.InProduction != srs.InProduction {
					slog.Debug("Updating in prod", "name", srs.Name, "mId", srs.MId, "old", srs.InProduction, "new", res.InProduction)
					err := srv.updateInProd(int(srs.ID), res.InProduction)
					if err != nil {
						slog.Warn("Failed to update in prod", "name", srs.Name, "mId", srs.MId, "old", srs.InProduction, "new", res.InProduction, "err", err)
					}
				}
			}()

			mu.Lock()
			fsd = append(fsd, res)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return fsd, nil
}

// Sorted by max(TvShowAdded, LastEpisodeAired,, LastEpisodeWatched)
func (srv *SeriesService) Feed() ([]*SeriesFeedItem, error) {

	series, err := srv.db.SeriesFeed()
	if err != nil {
		return nil, err
	}

	slog.Debug("Tv shows for feed", "series count", len(series))

	freshSeriesData, err := srv.freshSeriesData(series)

	slog.Debug("Fetched Data from Tmdb", "count", len(freshSeriesData))

	return srv.MakeFeedList(series, freshSeriesData), nil
}

func (srv *SeriesService) MakeFeedList(series []db.TvSeries, freshSeriesData []*SeriesFullItem) []*SeriesFeedItem {
	rv := []*SeriesFeedItem{}

	for _, srs := range series {

		idx := slices.IndexFunc(freshSeriesData, func(fd *SeriesFullItem) bool {
			return fd.Id == int(srs.MId)
		})

		fd := freshSeriesData[idx]
		watched := slices.CompactFunc(fd.EpsWatched, func(a, b SeriesEpisode) bool {
			return a.S == b.S && a.E == b.E
		})

		// update series data from fresh data
		srs.Name = fd.Name
		srs.Overview = fd.Overview
		srs.Year = fd.Year

		upNxt, _, _ := srv.MakeUpNext(int(srs.MId), fd)
		if upNxt == nil {
			continue
		}

		recentlyAired := false

		DaysAgo14 := time.Now().Add(-2 * time.Hour * 24 * 7)

		if fd.LastEpisodeToAir.AirDate.After(DaysAgo14) {
			recentlyAired = true
		}

		lastEpisodeWatchDate := time.Time{}
		if len(fd.EpsWatched) > 0 {
			lastEpisodeWatchDate = fd.EpsWatched[0].CreatedAt
		}

		rv = append(rv, &SeriesFeedItem{
			TvSeries:             srs,
			EpisodesTotal:        fd.NumberOfEpisodes,
			EpisodesAired:        fd.EpisodesAired,
			EpisodesWatched:      len(watched),
			UpNext:               upNxt,
			RecentlyAired:        recentlyAired,
			LastEpisodeAirDate:   fd.LastEpisodeToAir.AirDate,
			Image:                fd.ImgPoster,
			ShowAddDate:          srs.CreatedAt,
			LastEpisodeWatchDate: lastEpisodeWatchDate,
		})

		latest := func(a, b, c time.Time) time.Time {
			rv := a
			if b.After(rv) {
				rv = b
			}
			if c.After(rv) {
				rv = c
			}
			return rv
		}

		slices.SortFunc(rv, func(a, b *SeriesFeedItem) int {
			aT := latest(a.LastEpisodeAirDate, a.LastEpisodeWatchDate, a.ShowAddDate)
			bT := latest(b.LastEpisodeAirDate, b.LastEpisodeWatchDate, b.ShowAddDate)
			if aT.After(bT) {
				return -1
			} else if aT.Before(bT) {
				return 1
			}
			return 0
		})

	}

	return rv

}

func (srv *SeriesService) updateInProd(id int, inProd bool) error {
	return srv.db.SeriesUpdateInProd(id, inProd)
}

func (srv *SeriesService) UpNext(mId int) (*SeriesEpisode, error) {

	fd, err := srv.GetDetails(mId, false)
	if err != nil {
		slog.Error("Error while fetching data", "error", err, "mId", mId)
		return nil, err
	}

	rv, _, _ := srv.MakeUpNext(mId, fd)
	return rv, nil
}

// null if watched everything
func (srv *SeriesService) MakeUpNext(mId int, fd *SeriesFullItem) (*SeriesEpisode, int, int) {

	rv := SeriesEpisode{}

	var isWatched = func(s int, e int) bool {
		return slices.ContainsFunc(fd.EpsWatched, func(a SeriesEpisode) bool {
			return a.S == s && a.E == e
		})
	}

	upNextS := 1
	upNextE := 1

	lastWatchedFound := false

	lastAiredS := fd.LastEpisodeToAir.SeasonNumber
	lastAiredE := fd.LastEpisodeToAir.EpisodeNumber

	slices.SortFunc(fd.Seasons, func(a, b meta.TvSeason) int { return cmp.Compare(b.SeasonNumber, a.SeasonNumber) })

	for _, sn := range fd.Seasons {
		if !isLegitSeason(sn.Name, sn.SeasonNumber, fd.NumberOfSeasons) {
			continue
		}

		var eps int
		if sn.SeasonNumber < lastAiredS {
			eps = sn.EpisodeCount
		} else if sn.SeasonNumber == lastAiredS {
			eps = lastAiredE
		}

		for eNo := eps; eNo >= 1 && !lastWatchedFound; eNo -= 1 {
			if isWatched(sn.SeasonNumber, eNo) {
				lastWatchedFound = true
				break
			} else {
				upNextS = sn.SeasonNumber
				upNextE = eNo
			}
		}
	}

	if len(fd.EpsWatched) == fd.EpisodesAired || isWatched(lastAiredS, lastAiredE) {
		return nil, len(fd.EpsWatched), fd.EpisodesAired
	}

	rv.S = upNextS
	rv.E = upNextE

	return &rv, len(fd.EpsWatched), fd.EpisodesAired
}

// Returns insertModel from db
func (srv *SeriesService) Add(mId int, source string) (*db.TvSeries, error) {
	tvM, err := srv.GetTvMeta(mId)
	if err != nil {
		return nil, err
	}

	tvDb := &db.TvSeries{
		MName:          srv.meta.Name(),
		MId:            int64(tvM.Id),
		Name:           tvM.Name,
		Overview:       tvM.Overview,
		Year:           tvM.Year,
		TrackingStatus: db.TvStatusWatching,
		RuntimeApprox:  tvM.LastEpisodeToAir.Runtime,
		Source:         source,
		InProduction:   tvM.InProduction,
	}

	insertId, err := srv.db.SeriesAdd(tvDb)

	tvDb.ID = uint(insertId)

	return tvDb, err
}

func (srv *SeriesService) AddImportedEpisode(iep *ImportedTrackedEps) (int, error) {

	exists, err := srv.db.ImportedTrackedEpsCheck(iep.Key)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, ErrDuplicate
	}

	if iep.MId == 0 {
		return 0, errors.Join(ErrInvalidData, errors.New("Missing MId"))
	}

	ep, err := srv.getEpisodeDetails(iep.MId, iep.Season, iep.Episode)
	if err != nil {
		return -1, err
	}

	trackItem := db.TvTrackedEps{
		MName:      ep.MName,
		EpisodeMId: int64(ep.MId),
		SeriesMId:  int64(iep.MId),
		SeriesName: iep.SeriesName,
		Name:       ep.Name,
		Overview:   ep.Overview,
		Season:     ep.Season,
		Episode:    ep.Episode,
		Runtime:    ep.Runtime,
		Source:     db.SourceImport,
		SourceKey:  iep.Key,
	}

	trackItem.CreatedAt = iep.CreatedAt

	return srv.db.SeriesTrackedEpsAdd(&trackItem)

}

func (srv *SeriesService) UpdateStatus(mId int, newStatus db.TvStatus) error {

	if err := srv.db.SeriesStatusUpdate(mId, newStatus); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return errors.Join(ErrNotFound, db.ErrNotFound)
		}
		return err
	}
	return nil
}

func (srv *SeriesService) Remove(mId int) error {
	if err := srv.db.SeriesRem(mId); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return errors.Join(err, ErrNotFound)
		}
		return err
	}
	return nil
}

// Deleted Row Count
func (srv *SeriesService) SetEpisodeUnwatch(mId int, sNo int, eNo int) error {

	err := srv.db.SeriesTrackedEpRemove(mId, sNo, eNo)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			slog.Warn("Record Not found, but a delete request was issued", "mid", mId, "sno", sNo, "eno", eNo)
			return nil
		}
		return err
	}

	return nil

}

func (srv *SeriesService) SetEpisodeWatched(mId int, sNo int, eNo int, source string) (int, error) {

	status, err := srv.db.SeriesStatusGet(mId)

	if err != nil {
		return -1, err
	}

	if status == db.TvStatusNotWatching {
		slog.Debug("Tv show isn't added, creating a record for tracking", "mId", mId)
		_, err := srv.Add(mId, db.SourceUI)
		if err != nil {
			return -1, err
		}
	}

	ep, err := srv.getEpisodeDetails(mId, sNo, eNo)
	if err != nil {
		return -1, err
	}

	trackItem := db.TvTrackedEps{
		MName:      ep.MName,
		EpisodeMId: int64(ep.MId),
		SeriesMId:  int64(mId),
		SeriesName: ep.SeriesName,
		Name:       ep.Name,
		Overview:   ep.Overview,
		Season:     ep.Season,
		Episode:    ep.Episode,
		Runtime:    ep.Runtime,
		Source:     source,
	}

	return srv.db.SeriesTrackedEpsAdd(&trackItem)

}

// returns with all episode details in that season
func (srv *SeriesService) cacheSeasonInDb(mId int, season int, forceRefresh bool) (*db.TvSeason, error) {

	sn, err := srv.db.SeriesSeasonGet(mId, season)

	if err == nil && forceRefresh {
		err = db.ErrNotFound
	}

	if errors.Is(err, db.ErrNotFound) {
		slog.Debug("Caching Tv Season in Db", "mId", mId, "season", season)
		mSd, err := srv.meta.GetTVSeasonDetails(mId, season)
		if err != nil {
			return nil, err
		}

		srs, err := srv.GetTvMeta(mId)
		if err != nil {
			return nil, err
		}

		sn = MetaToDbSeason(srs, mSd)
		err = srv.db.SeriesSeasonAdd(sn)
	} else if err != nil {
		return nil, err
	}

	return sn, nil
}

func (srv *SeriesService) getEpisodeDetails(mId int, season int, episode int) (*db.TvEpisode, error) {

	ep, err := srv.db.SeriesEpisodeGet(mId, season, episode)

	forceRefresh := false
	if ep != nil {
		if ep.UpdatedAt.Before(ep.AirDate) && time.Now().After(ep.AirDate) {
			forceRefresh = true
			slog.Debug("getEpisodeDetails", "force refresh", forceRefresh, "id", mId, "season", season, "episode", episode)
		}

	}

	if err == nil {
		return ep, nil
	}

	if !errors.Is(err, db.ErrNotFound) {
		slog.Error("err while getting ep details", "id", mId, "season", season, "episode", episode)
		return nil, err
	}

	sn, err := srv.cacheSeasonInDb(mId, season, forceRefresh)

	if err != nil {
		slog.Error("Coudn't cache season", "season", season, "error", err)
		return nil, err
	}
	idx := slices.IndexFunc(sn.Episodes, func(x db.TvEpisode) bool {
		return (x.Episode == episode)
	})

	if idx == -1 {
		slog.Error("Episode not found in season", "mId", mId, "season", season, "episode", episode)
		return nil, ErrNotFound
	}

	return &sn.Episodes[idx], err
}

func (srv *SeriesService) deriveStatus(mId int, cur db.TvStatus) (db.TvStatus, error) {
	if cur == db.TvStatusStopped {
		return cur, nil
	}
	rv := cur
	fd, err := srv.GetDetails(mId, false)
	if err != nil {
		return -1, err
	}
	if fd.EpisodesAired == len(fd.EpsWatched) {
		if fd.InProduction {
			rv = db.TvStatusUpToDate
		} else {
			rv = db.TvStatusCompleted
		}
	}
	return rv, nil
}

func (srv *SeriesService) IptImportTvTimeData(zipPath string) error {

	state.Import.StageCnt = 2
	state.Import.Stage = 1
	if zipPath != "" {
		spc, epc, err := srv.ipt.DumpTvTimeGdprData(zipPath)
		if err != nil {
			return err
		}
		slog.Info("IMPORTED COUNT", "series_processed_count", spc, "episodes_processed_count", epc)
	}

	ids, err := srv.ipt.GetUnmatched()
	if err != nil {
		return err
	}

	state.Import.ProcessingEpsCntTotal = len(ids.Episodes)
	state.Import.ProcessingSrsCntTotal = len(ids.Series)
	state.Import.Stage = 2

	slog.Info("UNMATCHED COUNT", "series", len(ids.Series), "episodes", len(ids.Episodes))
	slog.Info("Processing", "series", len(ids.Series))

	for _, srs := range ids.Series {
		state.Import.ProcessingSrsCnt += 1

		exists, err := srv.db.ImportedSeriesCheck(srs.Key)
		if err != nil {
			return err
		}

		if exists {
			slog.Info("Already processed, skipping", "series", srs.Name, "tvTimeId", srs.TvTimeId)
			continue
		}

		ttd, err := srv.cGetTVFromTvTimeId(srs.TvTimeId)
		if err != nil {
			slog.Error("Error getting meta", "series", srs.Name, "tvTimeId", srs.TvTimeId, "error", err)
			err1 := srv.ipt.SetSeriesUnresolved(srs.Key, err.Error())
			if err1 != nil {
				slog.Error("Failed to set series unresolved", "error", err1)
				return err1
			}
			continue
		}

		tvd, err := srv.GetTvMeta(ttd.Id)
		if err != nil {
			slog.Error("Error getting tv details", "series", srs.Name, "tvTimeId", srs.TvTimeId, "error", err)
			err1 := srv.ipt.SetSeriesUnresolved(srs.Key, err.Error())
			if err1 != nil {
				slog.Error("Failed to set series unresolved", "error", err1)
				return err1
			}
			continue
		}

		ts := db.TvStatusWatching
		if srs.IsStopped {
			ts = db.TvStatusStopped
		}

		tvDb := &db.TvSeries{
			MName:          srv.meta.Name(),
			MId:            int64(tvd.Id),
			Name:           tvd.Name,
			Overview:       tvd.Overview,
			Year:           tvd.Year,
			TrackingStatus: ts,
			RuntimeApprox:  tvd.LastEpisodeToAir.Runtime,
			InProduction:   tvd.InProduction,
			Source:         db.SourceImport,
			SourceKey:      srs.Key,
		}
		tvDb.CreatedAt = srs.CreatedAt

		_, err = srv.db.SeriesAdd(tvDb)

		if err != nil {
			err1 := srv.ipt.SetSeriesUnresolved(srs.Key, err.Error())
			if err1 != nil {
				slog.Error("Failed to set series as unresolved", "error", err1)
				return err1
			}
			slog.Error("Error adding series to db", "series", srs.Name, "tvTimeId", srs.TvTimeId, "error", err)
		} else {

			err1 := srv.ipt.SetImported(srs.Key)
			if err1 != nil {
				slog.Warn("Failed to mark episode as imported", "error", err1)
			}
			slog.Info("IMPORT match success", "series", srs.Name, "tvTimeId", srs.TvTimeId, srv.meta.Name()+"Id", tvd.Id)

		}

	}

	insertEpisode := func(epd *db.TvEpisode, iep *ImportedTrackedEps) error {
		trackItem := db.TvTrackedEps{
			MName:      srv.meta.Name(),
			EpisodeMId: int64(epd.MId),
			SeriesMId:  int64(epd.SeriesMId),
			SeriesName: epd.SeriesName,
			Name:       epd.Name,
			Overview:   epd.Overview,
			Season:     epd.Season,
			Episode:    epd.Episode,
			Runtime:    epd.Runtime,
			Source:     db.SourceImport,
			SourceKey:  iep.Key,
		}
		trackItem.CreatedAt = iep.CreatedAt

		_, err := srv.db.SeriesTrackedEpsAdd(&trackItem)

		if err == nil {
			slog.Info("IMPORT episode match success", "series", iep.SeriesName, "TvTimeId", iep.TvTimeId, "TvTimeEId", iep.TvTimeEId, "MId", epd.MId, "SeriesMId", epd.SeriesMId)
			err1 := srv.ipt.SetImported(iep.Key)
			if err1 != nil {
				slog.Warn("Failed to mark episode as imported", "error", err1)
			}
		} else {
			slog.Error("Error inserting episode to db", "series", iep.SeriesName, "tvTimeId", iep.TvTimeId, "tvTimeEId", iep.TvTimeEId, "error", err)
			err1 := srv.ipt.SetEpisodeUnresolved(iep.Key, err.Error())
			if err1 != nil {
				slog.Error("Fatal - Failed to set episode unresolved", "error", err1)
				return err1
			}
		}

		return nil
	}

	for _, iep := range ids.Episodes {
		state.Import.ProcessingEpsCnt += 1

		exists, err := srv.db.ImportedTrackedEpsCheck(iep.Key)
		if err != nil {
			slog.Error("Fatal - Failed to check if imported episode exists in app db", "error", err)
			return err
		}
		if exists {
			slog.Warn("Already processed, skipping", "episode", iep.SeriesName, "seriesTvTimeId", iep.TvTimeId, "epTvTimeEId", iep.TvTimeEId)
			continue
		}

		var epd *db.TvEpisode
		tvm, err := srv.cGetTVFromTvTimeId(iep.TvTimeId)
		if err != nil {
			slog.Error("Error getting series meta for episode match", "series", iep.SeriesName, "tvTimeId", iep.TvTimeId, "error", err)
		} else {
			epd, err = srv.getEpisodeDetails(tvm.Id, iep.Season, iep.Episode)
			if err != nil {
				slog.Warn("Episode not found with std epNo & snNo", "series", iep.SeriesName, "season", iep.Season, "episode", iep.Episode, "err", err)
			}
		}

		if epd != nil {
			err1 := insertEpisode(epd, iep)
			if err1 != nil {
				return err1
			}
			// a match was found regardless of insert success, we can continue
			continue
		}

		slog.Warn("Trying to get episode details using external id", "series", iep.SeriesName, "TvTimeEId", iep.TvTimeEId, "season", iep.Season, "episode", iep.Episode)

		epm, err := srv.cGetEpisodeFromTvTimeId(iep.TvTimeEId)

		if err != nil {
			slog.Error("Error getting meta", "series", iep.SeriesName, "tvTimeId", iep.TvTimeId, "TvTimeEId", iep.TvTimeEId, "error", err)
			err1 := srv.ipt.SetEpisodeUnresolved(iep.Key, err.Error())
			if err1 != nil {
				slog.Error("Fatal - Failed to set series unresolved", "error", err1)
				return err1
			}
			continue
		}

		epd = &db.TvEpisode{
			MName:      srv.meta.Name(),
			MId:        int64(epm.Id),
			SeriesMId:  int64(epm.ShowId),
			SeriesName: tvm.Name,
			Name:       epm.Name,
			Overview:   epm.Overview,
			Season:     epm.SeasonNumber,
			Episode:    epm.EpisodeNumber,
			Runtime:    epm.Runtime,
		}

		err = insertEpisode(epd, iep)
		if err != nil {
			return err
		}
	}

	return err
}

func (srv *SeriesService) cGetTVFromTvTimeId(tvTimeId int) (*meta.TvDetails, error) {
	what := "TvTimeId"
	key := strconv.Itoa(tvTimeId)
	cd, err := srv.db.CacheGet(what, key)

	var refresh = func() (*db.Cached, error) {
		m, err := srv.meta.GetTVFromTvTimeId(tvTimeId)
		if err != nil {
			return nil, err
		}
		jsonStr, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}

		rv := &db.Cached{
			What:     what,
			Key:      key,
			JsonData: string(jsonStr),
		}

		err = srv.db.CacheSet(rv)
		if err != nil {
			return nil, err
		}
		return rv, nil
	}

	if errors.Is(err, db.ErrNotFound) {
		cd, err = refresh()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	expiredAt := cd.UpdatedAt.Add(24 * 7 * time.Hour)

	if time.Now().After(expiredAt) {
		slog.Debug("Cache expired, refreshing", "expiredAt", expiredAt, "UpdatedAt", cd.UpdatedAt)
		cd, err = refresh()
		if err != nil {
			return nil, err
		}
	}

	rv := meta.TvDetails{}
	err = json.Unmarshal([]byte(cd.JsonData), &rv)
	if err != nil {
		return nil, err
	}

	return &rv, nil

}

func (srv *SeriesService) cGetEpisodeFromTvTimeId(tvTimeId int) (*meta.TvEpisode, error) {
	what := "TvTimeEId"
	key := strconv.Itoa(tvTimeId)
	cd, err := srv.db.CacheGet(what, key)

	var refresh = func() (*db.Cached, error) {
		m, err := srv.meta.GetEpisodeFromTvTimeId(tvTimeId)
		if err != nil {
			return nil, err
		}
		jsonStr, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}

		rv := &db.Cached{
			What:     what,
			Key:      key,
			JsonData: string(jsonStr),
		}

		err = srv.db.CacheSet(rv)
		if err != nil {
			return nil, err
		}
		return rv, nil
	}

	if errors.Is(err, db.ErrNotFound) {
		cd, err = refresh()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	expiredAt := cd.UpdatedAt.Add(24 * 7 * time.Hour)

	if time.Now().After(expiredAt) {
		slog.Debug("Cache expired, refreshing", "expiredAt", expiredAt, "UpdatedAt", cd.UpdatedAt)
		cd, err = refresh()
		if err != nil {
			return nil, err
		}
	}

	rv := meta.TvEpisode{}
	err = json.Unmarshal([]byte(cd.JsonData), &rv)
	if err != nil {
		return nil, err
	}

	return &rv, nil

}

func (srv *SeriesService) IptGetUnresolved() (*ImportedData, error) {
	return srv.ipt.GetUnresolved()
}

func (srv *SeriesService) Stats() (*db.Stats, error) {
	return srv.db.SeriesStats()
}

func (srv *SeriesService) StatsMyShows(limit int) ([]db.StatsShow, error) {
	return srv.db.SeriesStatsMyShows(limit)
}

func (srv *SeriesService) Upcoming() ([]*UpcomingItem, error) {
	series, err := srv.db.SeriesWatchingInProdAll()
	if err != nil {
		return nil, err
	}
	fsd, err := srv.freshSeriesData(series)
	if err != nil {
		return nil, err
	}

	var rv []*UpcomingItem
	now := time.Now()

	for _, srs := range fsd {

		lastAiredS := srs.LastEpisodeToAir.SeasonNumber
		for _, sn := range srs.Seasons {
			if sn.SeasonNumber >= lastAiredS {
				epStart := 1
				if sn.SeasonNumber == lastAiredS {
					epStart = 1 + srs.LastEpisodeToAir.EpisodeNumber
				}
				for epNo := epStart; epNo <= sn.EpisodeCount; epNo++ {
					ep, err := srv.getEpisodeDetails(srs.Id, sn.SeasonNumber, epNo)
					if err != nil {
						slog.Error("Upcoming : Error getting episode details", "name", srs.Name, "season", sn.SeasonNumber, "episode", epNo, "err", err)
						return nil, err
					}
					if ep.AirDate.Before(now) {
						continue
					}
					rv = append(rv, &UpcomingItem{
						SeriesName: srs.Name,
						Year:       srs.Year,
						ImgPoster:  srs.ImgPoster,
						Episode:    ep,
					})
				}
			}
		}
	}

	slices.SortFunc(rv, func(a, b *UpcomingItem) int {
		if a.Episode.AirDate.Equal(b.Episode.AirDate) {
			if a.SeriesName == b.SeriesName {
				return cmp.Compare(a.Episode.Episode, b.Episode.Episode)
			}
			return strings.Compare(a.SeriesName, b.SeriesName)
		}
		return a.Episode.AirDate.Compare(b.Episode.AirDate)
	})

	return rv, nil

}

func (srv *SeriesService) GetPoster(mId int) string {
	m, err := srv.GetTvMeta(mId)
	if err != nil {
		slog.Error("GetPoster", "err", err)
		return ""
	}
	return m.ImgPoster
}
