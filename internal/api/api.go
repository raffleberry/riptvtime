package api

import (
	"log/slog"
	"net/http"

	"gitlab.com/raffleberry/riptvtime/internal/db"
	"gitlab.com/raffleberry/riptvtime/internal/meta"
	"gitlab.com/raffleberry/riptvtime/internal/ui"
)

type Api struct {
	Router *http.ServeMux
	db     db.Db
	meta   meta.Meta
}

func NewApi(db db.Db, meta meta.Meta) *Api {
	mux := http.NewServeMux()

	a := &Api{
		Router: mux,
		db:     db,
		meta:   meta,
	}

	// mux.HandleFunc("GET /api/series", apiSeriesAll())
	mux.HandleFunc("GET /api/series/search", a.SeriesSearch())
	mux.HandleFunc("GET /api/series/feed", a.SeriesFeed())

	mux.HandleFunc("POST /api/series", a.SeriesAdd())
	mux.HandleFunc("DELETE /api/series/{mId}", a.SeriesRem())

	mux.HandleFunc("PUT /api/series/{mId}/status", a.SeriesUpdateStatus())

	mux.HandleFunc("POST /api/series/episode", a.SeriesEpisodeWatch())
	mux.HandleFunc("PUT /api/series/episode", a.SeriesEpisodeUnWatch())

	// mux.HandleFunc("GET /api/series/{tmdbId}", apiSeriesGet())
	// mux.HandleFunc("GET /api/series/{tmdbId}/{episode}", apiSeriesEpisode())

	mux.Handle("GET /", ui.NewSpaHandler("internal/ui/static"))
	return a
}

func (a *Api) CacheSeason(mId int64, season int) error {
	mSd, err := a.meta.GetTVSeasonDetails(int(mId), season)
	slog.Debug("Caching Tv Season", "name", mSd.Name, "MSource", a.meta.Name(), "MId", mId, "Season", season, "Episodes", len(mSd.Episodes))

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

	_, err = a.db.SeriesSeasonAdd(&dbSd)

	return err
}

func (a *Api) addNewSeries(mId int) (int, error) {
	tvM, err := a.meta.GetTvDetails(mId)
	if err != nil {
		return 0, err
	}

	tvDb := &db.TvSeries{
		MName:          a.meta.Name(),
		MId:            int64(tvM.Id),
		Name:           tvM.Name,
		Overview:       tvM.Overview,
		Year:           tvM.Year,
		TrackingStatus: db.TvStatusWatching,
	}

	return a.db.SeriesAdd(tvDb)
}

func (a *Api) getEpisodeDetails(id int, season int, episode int) (*db.TvEpisode, error) {

	exists, err := a.db.SeriesEpisodeExists(id, season, episode)

	if err != nil {
		return nil, err
	}

	if !exists {
		err := a.CacheSeason(int64(id), season)
		if err != nil {
			return nil, err
		}
	}

	return a.db.SeriesEpisodeGet(id, season, episode)
}
