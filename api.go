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

type TvSeriesFeedItem struct {
	db.TvSeries
	EpisodesTotal      int
	EpisodesAired      int
	EpisodesWatched    int
	NextEpisodeToWatch db.TvEpisode
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
			TmdbId int `json:"tmdbId"`
		}

		if err := json.NewDecoder(c.R.Body).Decode(&payload); err != nil {
			return err
		}

		slog.Debug("tv add", "payload", payload)

		td, err := db.TmdbClient.GetTVDetails(payload.TmdbId, nil)

		if err != nil {
			return err
		}

		ts := db.TvSeries{
			TmdbId:         td.ID,
			Name:           td.Name,
			Overview:       td.Overview,
			Genres:         db.GenreToStr(td.Genres),
			Year:           db.ParseYear(td.FirstAirDate),
			FirstAirDate:   db.ParseAirDate(td.FirstAirDate),
			TrackingStatus: db.TvStatusWatching,
		}

		res := db.Conn.Create(&ts)

		if res.Error != nil {
			return res.Error
		}

		return c.JSON(http.StatusOK, ts)
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

func apiSeriesFeed() http.HandlerFunc {
	return server.WithCtx(func(c *server.Context) error {
		var series []db.TvSeries
		res := db.Conn.Find(&series).Where("tracking_status = ?", db.TvStatusWatching)
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

		var respItem []*TvSeriesFeedItem

		for _, show := range series {

			slog.Debug("::::Start Calculating Resp data")

			var trackedEps []db.TvTracking
			res := db.Conn.Find(&trackedEps).Where("series_tmdb_id = ?", show.TmdbId)
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

			var nextEpisodeToWatch *db.TvEpisode

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
				if nextEpisodeToWatch == nil {
					for epNo := 1; epNo <= s.EpisodeCount; epNo++ {
						key := fmt.Sprintf("%d-%d", s.SeasonNumber, epNo)
						if _, seen := watched[key]; !seen {
							nextEpisodeToWatch = &db.TvEpisode{
								TmdbID:     -1,
								SeriesTmdb: show.TmdbId,
								Season:     s.SeasonNumber,
								Episode:    epNo,
							}
							slog.Debug("Not Watched, Selecting & Breaking", "series", s.Name, "s", s.SeasonNumber, "e", epNo, "nextEpisodeToWatch", *nextEpisodeToWatch)
							break
						}
						slog.Debug("Watched", "series", s.Name, "s", s.SeasonNumber, "e", epNo, "key", key)
					}
				}
			}

			if nextEpisodeToWatch == nil {
				nextEpisodeToWatch = &db.TvEpisode{
					TmdbID:  show.TmdbId,
					Season:  1,
					Episode: 1,
				}
			}

			respItem = append(respItem, &TvSeriesFeedItem{
				TvSeries:           show,
				EpisodesTotal:      fd.NumberOfEpisodes,
				EpisodesAired:      episodesAired,
				EpisodesWatched:    len(trackedEps),
				NextEpisodeToWatch: *nextEpisodeToWatch,
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
