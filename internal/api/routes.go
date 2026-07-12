package api

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"gitlab.com/raffleberry/riptvtime/internal/db"
	"gitlab.com/raffleberry/riptvtime/internal/meta"
	"golang.org/x/sync/errgroup"
)

var (
	ErrInvalidRequest = errors.New("Invalid Request")
)

type tvSeriesFeedItem struct {
	db.TvSeries
	EpisodesTotal   int
	EpisodesAired   int
	EpisodesWatched int
	UpNextS         int
	UpNextE         int
	RecentlyAired   bool
}

type tvSeriesSearchItem struct {
	meta.TvSearchResult
	Status db.TvStatus
}
type tvSeriesSearchResult struct {
	Page         int
	Results      []tvSeriesSearchItem
	TotalPages   int
	TotalResults int
}

// queryParams = { q: query, p: page }
func (a *Api) SeriesSearch() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		urlVals := c.R.URL.Query()
		query := urlVals.Get("q")
		page := urlVals.Get("p")
		slog.Debug("search", "page", page, "query", query)
		metaRes, err := a.meta.Search(query, page)
		if err != nil {
			return err
		}

		respRes := tvSeriesSearchResult{
			Page:         metaRes.Page,
			TotalPages:   metaRes.TotalPages,
			TotalResults: metaRes.TotalResults,
		}

		for _, r := range metaRes.Results {

			status, err := a.db.SeriesStatusGet(r.Id)
			if err != nil {
				return err
			}

			item := tvSeriesSearchItem{
				TvSearchResult: r,
				Status:         status,
			}

			respRes.Results = append(respRes.Results, item)

		}

		return c.JSON(http.StatusOK, respRes)
	})
}

func (a *Api) SeriesAdd() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		var payload struct {
			MId   int    `json:"m_id"`
			MName string `json:"m_name"`
		}

		if err := json.NewDecoder(c.R.Body).Decode(&payload); err != nil {
			return err
		}

		slog.Debug("tv add", "payload", payload)

		if payload.MName != a.meta.Name() {
			return c.Error(http.StatusServiceUnavailable, "Requested Metadata Service is Unavailable")
		}

		insertId, err := a.addNewSeries(payload.MId)

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, insertId)
	})
}

// func apiSeriesGet() http.HandlerFunc {
// 	return WithCtx(func(c *Context) error {
// 		tmdbIdStr := c.R.PathValue("tmdbId")
// 		slog.Debug("tv get", "tmdbId", tmdbIdStr)

// 		tmdbId, err := strconv.Atoi(tmdbIdStr)
// 		if err != nil {
// 			c.W.WriteHeader(http.StatusBadRequest)
// 			return err
// 		}

// 		res, err := TmdbClient.GetTVDetails(tmdbId, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return c.JSON(http.StatusOK, NewTvShowFromTmdb(res))
// 	})
// }

// func apiSeriesEpisode() http.HandlerFunc {
// 	return WithCtx(func(c *Context) error {
// 		tmdbId := c.R.PathValue("tmdbId")
// 		episode := c.R.PathValue("episode")
// 		return nil
// 	})
// }

func (a *Api) SeriesEpisodeWatch() http.HandlerFunc {
	return WithCtx(func(ctx *Context) error {
		var payload struct {
			SeriesMId int `json:"series_m_id"`
			SeasonNo  int `json:"season_no"`
			EpisodeNo int `json:"episode_no"`
		}

		if err := json.NewDecoder(ctx.R.Body).Decode(&payload); err != nil {
			return err
		}

		status, err := a.db.SeriesStatusGet(payload.SeriesMId)

		if err != nil {
			return err
		}

		if status == db.TvStatusNotWatching {
			slog.Debug("Tv show isn't added, creating a record for tracking", "mId", payload.SeriesMId)
			_, err := a.addNewSeries(payload.SeriesMId)
			if err != nil {
				return err
			}
		}

		ep, err := a.getEpisodeDetails(payload.SeriesMId, payload.SeasonNo, payload.EpisodeNo)
		if err != nil {
			return err
		}

		trackItem := db.TvTrackedEps{
			MName:      ep.MName,
			EpisodeMId: ep.MId,
			SeriesMId:  int64(payload.SeriesMId),
			Name:       ep.Name,
			Overview:   ep.Overview,
			Season:     payload.SeasonNo,
			Episode:    payload.EpisodeNo,
			Runtime:    ep.Runtime,
		}

		insertId, err := a.db.SeriesTrackedEpsAdd(&trackItem)

		if err != nil {
			slog.Error("Error while creating a TvTracking Entry", "err", err)
			return err
		}

		return ctx.JSON(http.StatusOK, insertId)
	})
}

func (a *Api) SeriesEpisodeUnWatch() http.HandlerFunc {
	return WithCtx(func(ctx *Context) error {
		var payload struct {
			SeriesMId int `json:"series_m_id"`
			SeasonNo  int `json:"season_no"`
			EpisodeNo int `json:"episode_no"`
		}

		if err := json.NewDecoder(ctx.R.Body).Decode(&payload); err != nil {
			return err
		}

		deletedCnt, err := a.db.SeriesTrackedEpRemove(payload.SeriesMId, payload.SeasonNo, payload.EpisodeNo)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, deletedCnt)
	})
}

func (a *Api) SeriesFeed() http.HandlerFunc {
	return WithCtx(func(c *Context) error {

		series, err := a.db.SeriesWatchingAll()
		if err != nil {
			return err
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

				res, err := a.meta.GetTvDetails(int(mId))
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
			return err
		}

		slog.Debug("Fetched Data from Tmdb", "count", len(freshSeriesData))

		var respItem []*tvSeriesFeedItem

		for _, srs := range *series {

			slog.Debug("::::Start Calculating Resp data", "Series Name", srs.Name)
			trackedEps, err := a.db.SeriesTrackedEps(int(srs.MId))
			if err != nil {
				return err
			}
			watched := make(map[string]struct{})
			for _, t := range *trackedEps {
				key := fmt.Sprintf("%d-%d", t.Season, t.Episode)
				watched[key] = struct{}{}
			}

			// TODO: making this loop work while fetching fresh data
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
				eps := 0
				if sn.SeasonNumber < lastAiredS {
					eps = sn.EpisodeCount
				} else if sn.SeasonNumber == lastAiredS {
					eps = lastAiredE
				}
				episodesAired += eps

				for eNo := eps; eNo >= 1 && !lastWatchedFound; eNo -= 1 {
					if _, ok := watched[fmt.Sprintf("%d-%d", sn.SeasonNumber, eNo)]; ok {
						lastWatchedFound = true
						break
					} else {
						upNextS = sn.SeasonNumber
						upNextE = eNo
					}
				}
			}

			if len(watched) == episodesAired {
				continue
			}

			recentlyAired := false

			DaysAgo14 := time.Now().Add(-2 * time.Hour * 24 * 7)

			if fd.LastEpisodeToAir.AirDate.After(DaysAgo14) {
				recentlyAired = true
			}

			respItem = append(respItem, &tvSeriesFeedItem{
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

		return c.JSON(http.StatusOK, respItem)
	})
}

func (a *Api) SeriesUpdateStatus() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		mIdStr := c.R.PathValue("mId")

		mId, err := strconv.Atoi(mIdStr)

		if err != nil {
			return c.Error(http.StatusBadRequest, err.Error())
		}

		var payload struct {
			Status int `json:"status"`
		}

		if err := json.NewDecoder(c.R.Body).Decode(&payload); err != nil {
			return err
		}

		if payload.Status < 0 || payload.Status > len(db.TvStatusVals) {
			return c.Error(http.StatusBadRequest, "Invalid Status")
		}

		status, err := a.db.SeriesStatusGet(mId)
		if err != nil {
			return err
		}

		if status == db.TvStatusNotWatching {
			return c.Error(http.StatusBadRequest, "Tv show isn't added, cannot update status")
		}

		return a.db.SeriesStatusUpdate(mId, db.TvStatus(payload.Status))
	})
}

func (a *Api) SeriesRem() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		mIdStr := c.R.PathValue("mId")

		mId, err := strconv.Atoi(mIdStr)

		if err != nil {
			return c.Error(http.StatusBadRequest, err.Error())
		}

		if err := a.db.SeriesRem(mId); err != nil {
			if err == db.ErrNotFound {
				return c.Error(http.StatusNotFound, err.Error())
			}
			return err
		}

		return nil
	})
}
