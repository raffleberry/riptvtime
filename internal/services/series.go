package services

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"gitlab.com/raffleberry/riptvtime/internal/db"
	"gitlab.com/raffleberry/riptvtime/internal/meta"
	"gitlab.com/raffleberry/riptvtime/internal/utils"
	"golang.org/x/sync/errgroup"
)

var (
	ErrInvalidData = errors.New("Invalid Data")
	ErrNotFound    = errors.New("Not Found")
)

type SeriesService struct {
	db   db.Db
	meta meta.Meta
	c    db.Cache
}

func NewTvService(c db.Cache, db db.Db, meta meta.Meta) *SeriesService {
	return &SeriesService{
		db:   db,
		meta: meta,
		c:    c,
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

// Get Tv Show Details (Freshiestest Data possible)
func (srv *SeriesService) GetDetails(mId int, withEps bool) (*SeriesFullItem, error) {
	key := fmt.Sprintf("TvDetails{MId:%d}", mId)
	cd, err := srv.c.Get(key)

	var refresh = func() (*db.Cached, error) {
		m, err := srv.meta.GetTvDetails(mId)
		if err != nil {
			return nil, err
		}
		jsonStr, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}

		rv := &db.Cached{
			Key:      key,
			JsonData: string(jsonStr),
		}

		err = srv.c.Set(rv)
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

	expiredAt := cd.UpdatedAt.Add(48 * time.Hour)

	if time.Now().After(expiredAt) {
		slog.Debug("Cache expired, refreshing", "expiredAt", expiredAt, "UpdatedAt", cd.UpdatedAt)
		cd, err = refresh()
		if err != nil {
			return nil, err
		}
	}

	res := meta.TvDetails{}
	err = json.Unmarshal([]byte(cd.JsonData), &res)
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

		if withEps {

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
					Year:          ep.AirDate.Year(),
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

	for _, tep := range *tEps {
		epsWatched = append(epsWatched, SeriesEpisode{
			S: tep.Season,
			E: tep.Episode,
		})
	}

	return &SeriesFullItem{
		&res,
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

func (srv *SeriesService) Feed() (*[]SeriesFeedItem, error) {

	series, err := srv.db.SeriesWatchingAll()
	if err != nil {
		return nil, err
	}

	slog.Debug("Tv shows in Db", "series count", len(*series))

	var freshSeriesData []*meta.TvDetails

	var mu sync.Mutex

	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(3)

	for _, s := range *series {
		mId := s.MId

		g.Go(func() error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			res, err := srv.GetDetails(int(mId), false)
			if err != nil {
				slog.Error("Error while fetching data", "error", err, "mId", mId)
				return err
			}

			mu.Lock()
			freshSeriesData = append(freshSeriesData, res.TvDetails)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	slog.Debug("Fetched Data from Tmdb", "count", len(freshSeriesData))

	var rv []SeriesFeedItem

	for _, srs := range *series {

		slog.Debug("::::Start Calculating Resp data", "Series Name", srs.Name)
		trackedEps, err := srv.db.SeriesTrackedEps(int(srs.MId))
		if err != nil {
			return nil, err
		}
		watched := make(map[string]struct{})
		for _, t := range *trackedEps {
			key := fmt.Sprintf("%d-%d", t.Season, t.Episode)
			watched[key] = struct{}{}
		}

		var isWatched = func(s int, e int) bool {
			key := fmt.Sprintf("%d-%d", s, e)
			_, ok := watched[key]
			return ok
		}

		slog.Debug("watched eps", "name", srs.Name, "watched map", watched)

		idx := slices.IndexFunc(freshSeriesData, func(fd *meta.TvDetails) bool {
			return fd.Id == int(srs.MId)
		})

		fd := freshSeriesData[idx]

		// update series data from fresh data
		srs.Name = fd.Name
		srs.Overview = fd.Overview
		srs.Year = fd.Year

		upNextS := 1
		upNextE := 1

		lastWatchedFound := false

		episodesAired := 0
		lastAiredS := fd.LastEpisodeToAir.SeasonNumber
		lastAiredE := fd.LastEpisodeToAir.EpisodeNumber

		slices.SortFunc(fd.Seasons, func(a, b meta.TvSeason) int { return cmp.Compare(b.SeasonNumber, a.SeasonNumber) })

		for _, sn := range fd.Seasons {
			if 1 > sn.SeasonNumber || sn.SeasonNumber > fd.NumberOfSeasons {
				continue
			}

			var eps int
			if sn.SeasonNumber < lastAiredS {
				eps = sn.EpisodeCount
			} else if sn.SeasonNumber == lastAiredS {
				eps = lastAiredE
			}
			episodesAired += eps

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

		if len(watched) == episodesAired || isWatched(lastAiredS, lastAiredE) {
			continue
		}

		recentlyAired := false

		DaysAgo14 := time.Now().Add(-2 * time.Hour * 24 * 7)

		if fd.LastEpisodeToAir.AirDate.After(DaysAgo14) {
			recentlyAired = true
		}

		rv = append(rv, SeriesFeedItem{
			TvSeries:           srs,
			EpisodesTotal:      fd.NumberOfEpisodes,
			EpisodesAired:      episodesAired,
			EpisodesWatched:    len(watched),
			UpNextS:            upNextS,
			UpNextE:            upNextE,
			RecentlyAired:      recentlyAired,
			LastEpisodeAirDate: fd.LastEpisodeToAir.AirDate,
		})

		slog.Debug("::::End Calculating Resp data", "Series Name", srs.Name)
	}

	slices.SortFunc(rv, func(a SeriesFeedItem, b SeriesFeedItem) int {
		return b.LastEpisodeAirDate.Compare(a.LastEpisodeAirDate)
	})

	return &rv, nil

}

func (srv *SeriesService) UpNext(mId int) (*SeriesEpisode, error) {

	fd, err := srv.GetDetails(mId, false)
	if err != nil {
		slog.Error("Error while fetching data", "error", err, "mId", mId)
		return nil, err
	}

	var rv SeriesEpisode

	watched := make(map[string]struct{})
	for _, t := range fd.EpsWatched {
		key := fmt.Sprintf("%d-%d", t.S, t.E)
		watched[key] = struct{}{}
	}

	var isWatched = func(s int, e int) bool {
		key := fmt.Sprintf("%d-%d", s, e)
		_, ok := watched[key]
		return ok
	}

	upNextS := 1
	upNextE := 1

	lastWatchedFound := false

	episodesAired := 0
	lastAiredS := fd.LastEpisodeToAir.SeasonNumber
	lastAiredE := fd.LastEpisodeToAir.EpisodeNumber
	slog.Debug("dbg", "fd", fd)
	slices.SortFunc(fd.Seasons, func(a, b meta.TvSeason) int { return cmp.Compare(b.SeasonNumber, a.SeasonNumber) })

	for _, sn := range fd.Seasons {
		if !isLegitSeason(sn.Name, sn.SeasonNumber, fd.NumberOfSeasons) {
			continue
		}
		slog.Debug("processing", "name", sn.Name, "sno", sn.SeasonNumber, "cnt", fd.NumberOfSeasons)

		var eps int
		if sn.SeasonNumber < lastAiredS {
			eps = sn.EpisodeCount
		} else if sn.SeasonNumber == lastAiredS {
			eps = lastAiredE
		}
		episodesAired += eps

		for eNo := eps; eNo >= 1 && !lastWatchedFound; eNo -= 1 {
			if isWatched(sn.SeasonNumber, eNo) {
				slog.Debug("found", "sno", sn.SeasonNumber, "eno", eNo)
				lastWatchedFound = true
				break
			} else {
				upNextS = sn.SeasonNumber
				upNextE = eNo
			}
		}
	}

	if len(watched) == episodesAired || isWatched(lastAiredS, lastAiredE) {
		return nil, ErrNotFound
	}

	rv.S = upNextS
	rv.E = upNextE

	return &rv, nil

}

// Returns insertId from db
func (srv *SeriesService) Add(mId int) (*db.TvSeries, error) {
	tvM, err := srv.GetDetails(mId, false)
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
	}

	_, err = srv.db.SeriesAdd(tvDb)

	return tvDb, err
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

func (srv *SeriesService) SetEpisodeWatched(mId int, sNo int, eNo int) (int, error) {

	status, err := srv.db.SeriesStatusGet(mId)

	if err != nil {
		return -1, err
	}

	if status == db.TvStatusNotWatching {
		slog.Debug("Tv show isn't added, creating a record for tracking", "mId", mId)
		_, err := srv.Add(mId)
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
		Name:       ep.Name,
		Overview:   ep.Overview,
		Season:     ep.Season,
		Episode:    ep.Episode,
		Runtime:    ep.Runtime,
	}

	return srv.db.SeriesTrackedEpsAdd(&trackItem)

}

// returns with all episode details in that season
func (srv *SeriesService) cacheSeasonInDb(mId int, season int) (*db.TvSeason, error) {

	sn, err := srv.db.SeriesSeasonGet(mId, season)

	if errors.Is(err, db.ErrNotFound) {
		mSd, err := srv.meta.GetTVSeasonDetails(mId, season)
		slog.Debug("Caching Tv Season in Db", "name", mSd.Name, "MSource", srv.meta.Name(), "MId", mId, "Season", mSd.SeasonNumber, "Episodes", len(mSd.Episodes))
		if err != nil {
			return nil, err
		}

		sn = MetaToDbSeason(mId, mSd)
		err = srv.db.SeriesSeasonAdd(sn)
	} else if err != nil {
		return nil, err
	}

	return sn, nil
}

func (srv *SeriesService) getEpisodeDetails(id int, season int, episode int) (*db.TvEpisode, error) {

	ep, err := srv.db.SeriesEpisodeGet(id, season, episode)

	if err == nil {
		return ep, nil
	}

	if !errors.Is(err, db.ErrNotFound) {
		return nil, errors.Join(err, errors.New(utils.Jn("err while getting ep details", id, season, episode)))
	}

	sn, err := srv.cacheSeasonInDb(id, season)

	if err != nil {
		return nil, errors.Join(err, errors.New(utils.Jn("Coudn't cache season", "season", season, "error", err)))
	}
	idx := slices.IndexFunc(sn.Episodes, func(x db.TvEpisode) bool {
		return (x.Episode == episode)
	})

	if idx == -1 {
		return nil, fmt.Errorf("%w: Episode %d Not Found in slice returned by db", ErrNotFound, episode)
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
