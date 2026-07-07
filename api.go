package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"sync"

	tmdb "github.com/cyruzin/golang-tmdb"
	"gitlab.com/raffleberry/riptvtime/db"
	"gitlab.com/raffleberry/riptvtime/server"
	"golang.org/x/sync/errgroup"
)

type tvSeriesFeedItem struct {
	db.TvSeries
	EpisodesTotal   int
	EpisodesAired   int
	EpisodesWatched int
	UpNextS         int
	UpNextE         int
	UpToDate        bool
}

func apiSeriesSearch() http.HandlerFunc {
	return server.WithCtx(func(c *server.Context) error {
		urlVals := c.R.URL.Query()
		query := urlVals.Get("q")
		page := urlVals.Get("p")
		slog.Debug("search", "page", page, "query", query)
		res, err := db.TmdbClient.GetSearchTVShow(query, map[string]string{"page": page})
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, res)
	})
}

// func apiSeriesUpdateStatus() http.HandlerFunc {
// 	return server.WithCtx(func(c *server.Context) error {
// 		tmdbId := c.R.PathValue("tmdbId")
// 		return nil
// 	})
// }

// func apiSeriesRemove() http.HandlerFunc {
// 	return server.WithCtx(func(c *server.Context) error {
// 		tmdbId := c.R.PathValue("tmdbId")
// 		return nil
// 	})
// }

func apiSeriesAdd() http.HandlerFunc {
	return server.WithCtx(func(c *server.Context) error {
		var payload struct {
			TmdbId int `json:"tmdb_id"`
		}

		if err := json.NewDecoder(c.R.Body).Decode(&payload); err != nil {
			return err
		}

		slog.Debug("tv add", "payload", payload)

		insertId, err := db.NewTrackTv(payload.TmdbId)

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, insertId)
	})
}

// func apiSeriesGet() http.HandlerFunc {
// 	return server.WithCtx(func(c *server.Context) error {
// 		tmdbIdStr := c.R.PathValue("tmdbId")
// 		slog.Debug("tv get", "tmdbId", tmdbIdStr)

// 		tmdbId, err := strconv.Atoi(tmdbIdStr)
// 		if err != nil {
// 			c.W.WriteHeader(http.StatusBadRequest)
// 			return err
// 		}

// 		res, err := db.TmdbClient.GetTVDetails(tmdbId, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return c.JSON(http.StatusOK, db.NewTvShowFromTmdb(res))
// 	})
// }

// func apiSeriesEpisode() http.HandlerFunc {
// 	return server.WithCtx(func(c *server.Context) error {
// 		tmdbId := c.R.PathValue("tmdbId")
// 		episode := c.R.PathValue("episode")
// 		return nil
// 	})
// }

func apiSeriesEpisodeWatch() http.HandlerFunc {
	return server.WithCtx(func(ctx *server.Context) error {
		var payload struct {
			SeriesTmdbId int `json:"series_tmdb_id"`
			SeasonNo     int `json:"season_no"`
			EpisodeNo    int `json:"episode_no"`
		}

		if err := json.NewDecoder(ctx.R.Body).Decode(&payload); err != nil {
			return err
		}

		var seriesLen int64
		res := db.Db.Model(&db.TvSeries{}).Where("tmdb_id = ?", payload.SeriesTmdbId).Count(&seriesLen)

		if res.Error != nil {
			return res.Error
		}

		if seriesLen == 0 {
			slog.Debug("Tv show isn't added, creating a record for tracking", "tmdb_id", payload.SeriesTmdbId)
			_, err := db.NewTrackTv(payload.SeriesTmdbId)
			if err != nil {
				return err
			}
		}

		ep, err := db.GetEpisodeDetails(int64(payload.SeriesTmdbId), payload.SeasonNo, payload.EpisodeNo)
		if err != nil {
			return err
		}

		trackItem := db.TvTracking{
			EpisodeTmdbId: ep.TmdbID,
			SeriesTmdbId:  ep.SeriesTmdb,
			Name:          ep.Name,
			Overview:      ep.Overview,
			Season:        ep.Season,
			Episode:       ep.Episode,
			Runtime:       ep.Runtime,
		}

		res = db.Db.Create(&trackItem)

		if res.Error != nil {
			slog.Error("Error while creating a TvTracking Entry", "err", res.Error)
			return res.Error
		}

		return ctx.JSON(http.StatusOK, trackItem)
	})
}

func apiSeriesEpisodeUnWatch() http.HandlerFunc {
	return server.WithCtx(func(ctx *server.Context) error {
		var payload struct {
			SeriesTmdbId int `json:"series_tmdb_id"`
			SeasonNo     int `json:"season_no"`
			EpisodeNo    int `json:"episode_no"`
		}

		if err := json.NewDecoder(ctx.R.Body).Decode(&payload); err != nil {
			return err
		}

		var isTracking bool

		err := db.Db.Model(&db.TvTracking{}).
			Select("count(*) > 0").
			Where("series_tmdb_id = ? AND season = ? AND episode = ?", payload.SeriesTmdbId, payload.SeasonNo, payload.EpisodeNo).
			Limit(1).
			Find(&isTracking).Error

		if err != nil {
			slog.Error("Failed to check if series is being tracked", "err", err)
			return err
		}

		if !isTracking {
			slog.Debug("Tv show isn't added, creating a record for tracking", "tmdb_id", payload.SeriesTmdbId)
			_, err := db.NewTrackTv(payload.SeriesTmdbId)
			if err != nil {
				return err
			}
		}

		var epTrack db.TvTracking
		err = db.Db.First(&epTrack, "series_tmdb_id = ? AND season = ? AND episode = ?", payload.SeriesTmdbId, payload.SeasonNo, payload.EpisodeNo).Error
		if err != nil {
			slog.Error("Failed to get episode tracking details", "err", err)
			return err
		}

		res := db.Db.Delete(&epTrack)

		if res.Error != nil {
			slog.Error("Error while deleting a TvTracking Entry", "err", res.Error)
			return res.Error
		}

		return ctx.JSON(http.StatusOK, struct {
			DeletedId int
		}{
			DeletedId: int(epTrack.ID),
		})
	})
}

func apiSeriesFeed() http.HandlerFunc {
	return server.WithCtx(func(c *server.Context) error {
		var series []db.TvSeries
		res := db.Db.Find(&series).Where("tracking_status = ?", db.TvStatusWatching)
		if res.Error != nil {
			return res.Error
		}

		slog.Debug("Tv shows in Db", "series count", len(series))

		var freshSeriesData []*tmdb.TVDetails

		var mu sync.Mutex

		g, ctx := errgroup.WithContext(context.Background())
		g.SetLimit(3)

		for _, s := range series {
			tmdbId := s.TmdbId

			g.Go(func() error {
				if ctx.Err() != nil {
					return ctx.Err()
				}

				res, err := db.TmdbClient.GetTVDetails(int(tmdbId), nil)
				if err != nil {
					slog.Error("Error while fetching data", "error", err, "tmdbId", tmdbId)
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

		for _, show := range series {

			slog.Debug("::::Start Calculating Resp data")

			var trackedEps []db.TvTracking
			res := db.Db.Find(&trackedEps).Where("series_tmdb_id = ?", show.TmdbId)
			watched := make(map[string]struct{})
			for _, t := range trackedEps {
				key := fmt.Sprintf("%d-%d", t.Season, t.Episode)
				watched[key] = struct{}{}
			}

			slog.Debug("Watched Episodes", "mp", watched)

			if res.Error != nil {
				return res.Error
			}
			idx := slices.IndexFunc(freshSeriesData, func(fd *tmdb.TVDetails) bool {
				return fd.ID == show.TmdbId
			})
			fd := freshSeriesData[idx]

			// update series data from fresh data
			show.FirstAirDate = db.ParseAirDate(fd.FirstAirDate)
			show.Genres = db.GenreToStr(fd.Genres)
			show.Name = fd.Name
			show.Overview = fd.Overview
			show.Year = db.ParseYear(fd.FirstAirDate)

			upNextS := 1
			upNextE := 1

			episodesAired := 0
			lastAiredS := fd.LastEpisodeToAir.SeasonNumber
			lastAiredE := fd.LastEpisodeToAir.EpisodeNumber

			for _, s := range fd.Seasons {
				if !strings.HasPrefix(strings.ToLower(s.Name), "season") {
					continue
				}
				if s.SeasonNumber < lastAiredS {
					episodesAired += s.EpisodeCount
				} else if s.SeasonNumber == lastAiredS {
					episodesAired += lastAiredE
				}
				for epNo := 1; epNo <= s.EpisodeCount; epNo++ {
					key := fmt.Sprintf("%d-%d", s.SeasonNumber, epNo)
					if _, seen := watched[key]; !seen {
						if s.SeasonNumber > upNextS || (s.SeasonNumber == upNextS && epNo > upNextE) {
							upNextS = s.SeasonNumber
							upNextE = epNo
						}
					}
				}
			}

			respItem = append(respItem, &tvSeriesFeedItem{
				TvSeries:        show,
				EpisodesTotal:   fd.NumberOfEpisodes,
				EpisodesAired:   episodesAired,
				EpisodesWatched: len(watched),
				UpNextS:         upNextS,
				UpNextE:         upNextE,
				UpToDate:        episodesAired == len(watched),
			})

			slog.Debug("::::End Calculating Resp data")
		}

		return c.JSON(http.StatusOK, respItem)
	})
}

// func apiSeriesAll() http.HandlerFunc {
// 	return server.WithCtx(func(c *server.Context) error {
// 		return nil
// 	})
// }
