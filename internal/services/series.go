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

// Get Tv Show Details
func (srv *SeriesService) GetDetails(mId int) (*meta.TvDetails, error) {
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

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			cd, err = refresh()
			if err != nil {
				return nil, err
			}
		}
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
	return &res, nil
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

			res, err := srv.GetDetails(int(mId))
			if err != nil {
				slog.Error("Error while fetching data", "error", err, "mId", mId)
				return err
			}

			mu.Lock()
			freshSeriesData = append(freshSeriesData, res)
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

		// TODO: making this loop work while fetching fresh data
		// Rethink after caching(maybe not required)

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
			TvSeries:        srs,
			EpisodesTotal:   fd.NumberOfEpisodes,
			EpisodesAired:   episodesAired,
			EpisodesWatched: len(watched),
			UpNextS:         upNextS,
			UpNextE:         upNextE,
			RecentlyAired:   recentlyAired,
		})

		slog.Debug("::::End Calculating Resp data", "Series Name", srs.Name)
	}

	return &rv, nil

}

// Returns insertId from db
func (srv *SeriesService) Add(mId int) (int, error) {
	tvM, err := srv.meta.GetTvDetails(mId)
	if err != nil {
		return 0, err
	}

	tvDb := &db.TvSeries{
		MName:          srv.meta.Name(),
		MId:            int64(tvM.Id),
		Name:           tvM.Name,
		Overview:       tvM.Overview,
		Year:           tvM.Year,
		TrackingStatus: db.TvStatusWatching,
	}

	return srv.db.SeriesAdd(tvDb)
}

func (srv *SeriesService) GetStatus(mId int) (db.TvStatus, error) {
	return srv.db.SeriesStatusGet(mId)
}

func (srv *SeriesService) UpdateStatus(mId int, newStatus db.TvStatus) error {

	if err := srv.db.SeriesStatusUpdate(mId, newStatus); err != nil {
		if err == db.ErrNotFound {
			return fmt.Errorf("Series with MId %d Not Found, %w", mId, errors.Join(ErrNotFound, db.ErrNotFound))
		}
		return err
	}
	return nil
}

func (srv *SeriesService) Remove(mId int) error {
	if err := srv.db.SeriesRem(mId); err != nil {
		if err == db.ErrNotFound {
			return fmt.Errorf("Series with MId %d Not Found, %w", mId, errors.Join(ErrNotFound, db.ErrNotFound))
		}
		return err
	}
	return nil
}

// Deleted Row Count
func (srv *SeriesService) SetEpisodeUnwatch(mId int, sNo int, eNo int) (int, error) {

	deletedCnt, err := srv.db.SeriesTrackedEpRemove(mId, sNo, eNo)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	return deletedCnt, nil

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
		EpisodeMId: ep.MId,
		SeriesMId:  ep.SeriesMId,
		Name:       ep.Name,
		Overview:   ep.Overview,
		Season:     ep.Season,
		Episode:    ep.Episode,
		Runtime:    ep.Runtime,
	}

	return srv.db.SeriesTrackedEpsAdd(&trackItem)

}

func (srv *SeriesService) cacheSeason(mId int64, season int) error {
	mSd, err := srv.meta.GetTVSeasonDetails(int(mId), season)
	slog.Debug("Caching Tv Season", "name", mSd.Name, "MSource", srv.meta.Name(), "MId", mId, "Season", season, "Episodes", len(mSd.Episodes))

	if err != nil {
		return err
	}

	dbEps := []db.TvEpisode{}
	for _, e := range mSd.Episodes {
		dbEps = append(dbEps, db.TvEpisode{
			MId:       int64(e.Id),
			MName:     e.MName,
			Name:      e.Name,
			SeriesMId: mId,
			Overview:  e.Overview,
			Season:    e.SeasonNumber,
			Episode:   e.EpisodeNumber,
			Runtime:   e.Runtime,
			AirDate:   e.AirDate,
		})
	}

	dbSd := db.TvSeason{
		MId:       int64(mSd.Id),
		MName:     mSd.MName,
		AirDate:   mSd.AirDate,
		SeriesMId: mId,
		Season:    season,
		Name:      mSd.Name,
		Overview:  mSd.Overview,
		Episodes:  dbEps,
	}

	_, err = srv.db.SeriesSeasonAdd(&dbSd)

	return err
}

func (srv *SeriesService) getEpisodeDetails(id int, season int, episode int) (*db.TvEpisode, error) {

	ep, err := srv.db.SeriesEpisodeGet(id, season, episode)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {

			err = srv.cacheSeason(int64(id), season)
			if err != nil {
				return nil, err
			}

			return srv.db.SeriesEpisodeGet(id, season, episode)
		} else {
			return nil, err
		}
	}

	return ep, err
}
