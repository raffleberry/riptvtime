package api

import (
	"errors"
	"log/slog"
	"net/http"

	"gitlab.com/raffleberry/riptvtime/internal/db"
	"gitlab.com/raffleberry/riptvtime/internal/meta"
	"gitlab.com/raffleberry/riptvtime/internal/ui"
)

var (
	ErrInvalidRequest = errors.New("Invalid Request")
)

type Api struct {
	Router *http.ServeMux
	db     *db.Db
	meta   meta.Meta
}

func NewApi(db *db.Db, meta meta.Meta) *Api {
	mux := http.NewServeMux()

	a := &Api{
		Router: mux,
		db:     db,
		meta:   meta,
	}

	// mux.HandleFunc("GET /api/series", apiSeriesAll())
	mux.HandleFunc("GET /api/series/search", a.SeriesSearch())
	// mux.HandleFunc("GET /api/series/feed", apiSeriesFeed())
	// mux.HandleFunc("GET /api/series/{tmdbId}", apiSeriesGet())
	// mux.HandleFunc("GET /api/series/{tmdbId}/{episode}", apiSeriesEpisode())
	// mux.HandleFunc("POST /api/series", apiSeriesAdd())
	// mux.HandleFunc("POST /api/series/episode", apiSeriesEpisodeWatch())
	// mux.HandleFunc("DELETE /api/series/episode", apiSeriesEpisodeUnWatch())
	// mux.HandleFunc("DELETE /api/series/{tmdbId}", apiSeriesRemove())
	// mux.HandleFunc("PUT /api/series/{tmdbId}", apiSeriesUpdateStatus())

	mux.Handle("GET /", ui.NewSpaHandler())

	return a
}

type tvSeriesFeedItem struct {
	db.TvSeries
	EpisodesTotal   int
	EpisodesAired   int
	EpisodesWatched int
	UpNextS         int
	UpNextE         int
	RecentlyAired   bool
}

// queryParams = { q: query, p: page }
func (a *Api) SeriesSearch() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		urlVals := c.R.URL.Query()
		query := urlVals.Get("q")
		page := urlVals.Get("p")
		slog.Debug("search", "page", page, "query", query)
		res, err := a.meta.Search(query, page)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, res)
	})
}

// func apiSeriesUpdateStatus() http.HandlerFunc {
// 	return WithCtx(func(c *Context) error {
// 		tmdbId := c.R.PathValue("tmdbId")
// 		return nil
// 	})
// }

// func apiSeriesRemove() http.HandlerFunc {
// 	return WithCtx(func(c *Context) error {
// 		tmdbId := c.R.PathValue("tmdbId")
// 		return nil
// 	})
// }

// func apiSeriesAdd() http.HandlerFunc {
// 	return WithCtx(func(c *Context) error {
// 		var payload struct {
// 			TmdbId int `json:"tmdb_id"`
// 		}

// 		if err := json.NewDecoder(c.R.Body).Decode(&payload); err != nil {
// 			return err
// 		}

// 		slog.Debug("tv add", "payload", payload)

// 		insertId, err := NewTrackTv(payload.TmdbId)

// 		if err != nil {
// 			return err
// 		}

// 		return c.JSON(http.StatusOK, insertId)
// 	})
// }

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

// func apiSeriesEpisodeWatch() http.HandlerFunc {
// 	return WithCtx(func(ctx *Context) error {
// 		var payload struct {
// 			SeriesTmdbId int `json:"series_tmdb_id"`
// 			SeasonNo     int `json:"season_no"`
// 			EpisodeNo    int `json:"episode_no"`
// 		}

// 		if err := json.NewDecoder(ctx.R.Body).Decode(&payload); err != nil {
// 			return err
// 		}

// 		var isAdded bool

// 		err := Db.Model(&TvSeason{}).
// 			Select("count(*) > 0").
// 			Where("tmdb_id = ?", payload.SeriesTmdbId).
// 			Limit(1).
// 			Find(&isAdded).Error

// 		if err != nil {
// 			slog.Error("Failed to check if tv series was added", "err", err)
// 			return err
// 		}

// 		if !isAdded {
// 			slog.Debug("Tv show isn't added, creating a record for tracking", "tmdb_id", payload.SeriesTmdbId)
// 			_, err := NewTrackTv(payload.SeriesTmdbId)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		ep, err := GetEpisodeDetails(int64(payload.SeriesTmdbId), payload.SeasonNo, payload.EpisodeNo)
// 		if err != nil {
// 			return err
// 		}

// 		trackItem := TvTracking{
// 			EpisodeTmdbId: ep.TmdbID,
// 			SeriesTmdbId:  ep.SeriesTmdb,
// 			Name:          ep.Name,
// 			Overview:      ep.Overview,
// 			Season:        ep.Season,
// 			Episode:       ep.Episode,
// 			Runtime:       ep.Runtime,
// 		}

// 		err = Db.Create(&trackItem).Error

// 		if err != nil {
// 			slog.Error("Error while creating a TvTracking Entry", "err", err)
// 			return err
// 		}

// 		return ctx.JSON(http.StatusOK, trackItem)
// 	})
// }

// func apiSeriesEpisodeUnWatch() http.HandlerFunc {
// 	return WithCtx(func(ctx *Context) error {
// 		var payload struct {
// 			SeriesTmdbId int `json:"series_tmdb_id"`
// 			SeasonNo     int `json:"season_no"`
// 			EpisodeNo    int `json:"episode_no"`
// 		}

// 		if err := json.NewDecoder(ctx.R.Body).Decode(&payload); err != nil {
// 			return err
// 		}

// 		var isTracking bool

// 		err := Db.Model(&TvTracking{}).
// 			Select("count(*) > 0").
// 			Where("series_tmdb_id = ? AND season = ? AND episode = ?", payload.SeriesTmdbId, payload.SeasonNo, payload.EpisodeNo).
// 			Limit(1).
// 			Find(&isTracking).Error

// 		if err != nil {
// 			slog.Error("Failed to check if episode was tracked", "err", err)
// 			return err
// 		}

// 		if !isTracking {
// 			return ctx.Error(http.StatusBadRequest, "Not Tracked")
// 		}

// 		var epTrack TvTracking
// 		err = Db.First(&epTrack, "series_tmdb_id = ? AND season = ? AND episode = ?", payload.SeriesTmdbId, payload.SeasonNo, payload.EpisodeNo).Error
// 		if err != nil {
// 			slog.Error("Failed to get episode tracking details", "err", err)
// 			return err
// 		}

// 		res := Db.Delete(&epTrack)

// 		if res.Error != nil {
// 			slog.Error("Error while deleting a TvTracking Entry", "err", res.Error)
// 			return res.Error
// 		}

// 		return ctx.JSON(http.StatusOK, epTrack.ID)
// 	})
// }

// func apiSeriesFeed() http.HandlerFunc {
// 	return WithCtx(func(c *Context) error {
// 		var series []TvSeries
// 		res := Db.Find(&series).Where("tracking_status = ?", TvStatusWatching)
// 		if res.Error != nil {
// 			return res.Error
// 		}

// 		slog.Debug("Tv shows in Db", "series count", len(series))

// 		var freshSeriesData []*tmdb.TVDetails

// 		var mu sync.Mutex

// 		g, ctx := errgroup.WithContext(context.Background())
// 		g.SetLimit(3)

// 		for _, s := range series {
// 			tmdbId := s.TmdbId

// 			g.Go(func() error {
// 				if ctx.Err() != nil {
// 					return ctx.Err()
// 				}

// 				res, err := TmdbClient.GetTVDetails(int(tmdbId), nil)
// 				if err != nil {
// 					slog.Error("Error while fetching data", "error", err, "tmdbId", tmdbId)
// 					return err
// 				}

// 				mu.Lock()
// 				freshSeriesData = append(freshSeriesData, res)
// 				mu.Unlock()

// 				return nil
// 			})
// 		}

// 		if err := g.Wait(); err != nil {
// 			return err
// 		}

// 		slog.Debug("Fetched Data from Tmdb", "count", len(freshSeriesData))

// 		var respItem []*tvSeriesFeedItem

// 		for _, show := range series {

// 			slog.Debug("::::Start Calculating Resp data", "Series Name", show.Name)

// 			var trackedEps []TvTracking
// 			res := Db.Find(&trackedEps).Where("series_tmdb_id = ?", show.TmdbId)
// 			watched := make(map[string]struct{})
// 			for _, t := range trackedEps {
// 				key := fmt.Sprintf("%d-%d", t.Season, t.Episode)
// 				watched[key] = struct{}{}
// 			}

// 			slog.Debug("Watched Episodes", "mp", watched)

// 			if res.Error != nil {
// 				return res.Error
// 			}
// 			idx := slices.IndexFunc(freshSeriesData, func(fd *tmdb.TVDetails) bool {
// 				return fd.ID == show.TmdbId
// 			})
// 			fd := freshSeriesData[idx]

// 			// update series data from fresh data
// 			show.FirstAirDate = ParseAirDate(fd.FirstAirDate)
// 			show.Genres = genreToStr(fd.Genres)
// 			show.Name = fd.Name
// 			show.Overview = fd.Overview
// 			show.Year = ParseYear(fd.FirstAirDate)

// 			upNextS := 1
// 			upNextE := 1

// 			episodesAired := 0
// 			lastAiredS := fd.LastEpisodeToAir.SeasonNumber
// 			lastAiredE := fd.LastEpisodeToAir.EpisodeNumber

// 			for _, s := range fd.Seasons {
// 				if !strings.HasPrefix(strings.ToLower(s.Name), "season") {
// 					continue
// 				}
// 				if s.SeasonNumber < lastAiredS {
// 					episodesAired += s.EpisodeCount
// 				} else if s.SeasonNumber == lastAiredS {
// 					episodesAired += lastAiredE
// 				}
// 				for epNo := 1; epNo <= s.EpisodeCount; epNo++ {
// 					key := fmt.Sprintf("%d-%d", s.SeasonNumber, epNo)
// 					if _, seen := watched[key]; !seen {
// 						if s.SeasonNumber > upNextS || (s.SeasonNumber == upNextS && epNo > upNextE) {
// 							upNextS = s.SeasonNumber
// 							upNextE = epNo
// 						}
// 					}
// 				}
// 			}

// 			if len(watched) == episodesAired {
// 				continue
// 			}

// 			recentlyAired := false

// 			DaysAgo14 := time.Now().Add(-2 * time.Hour * 24 * 7)
// 			if ParseAirDate(fd.LastEpisodeToAir.AirDate).After(DaysAgo14) {
// 				recentlyAired = true
// 			}

// 			respItem = append(respItem, &tvSeriesFeedItem{
// 				TvSeries:        show,
// 				EpisodesTotal:   fd.NumberOfEpisodes,
// 				EpisodesAired:   episodesAired,
// 				EpisodesWatched: len(watched),
// 				UpNextS:         upNextS,
// 				UpNextE:         upNextE,
// 				RecentlyAired:   recentlyAired,
// 			})

// 			slog.Debug("::::End Calculating Resp data", "Series Name", show.Name)
// 		}

// 		return c.JSON(http.StatusOK, respItem)
// 	})
// }

// func apiSeriesAll() http.HandlerFunc {
// 	return WithCtx(func(c *Context) error {
// 		return nil
// 	})
// }
