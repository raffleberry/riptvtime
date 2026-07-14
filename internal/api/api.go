package api

import (
	"net/http"

	"gitlab.com/raffleberry/riptvtime/internal/db"
	"gitlab.com/raffleberry/riptvtime/internal/meta"
	"gitlab.com/raffleberry/riptvtime/internal/services"
	"gitlab.com/raffleberry/riptvtime/internal/ui"
)

// if response is not 200, error should be sent in plain text as response
// with the appropriate error code as status
type Api struct {
	Router *http.ServeMux
	db     db.Db
	meta   meta.Meta
	tv     *services.SeriesService
}

func NewApi(db db.Db, meta meta.Meta, tv *services.SeriesService) *Api {
	mux := http.NewServeMux()

	a := &Api{
		Router: mux,
		db:     db,
		meta:   meta,
		tv:     tv,
	}

	// mux.HandleFunc("GET /api/series", apiSeriesAll())
	mux.HandleFunc("GET /api/series/search", a.SeriesSearch())
	mux.HandleFunc("GET /api/series/feed", a.SeriesFeed())

	mux.HandleFunc("GET /api/series/{mId}", a.SeriesGet())
	mux.HandleFunc("POST /api/series", a.SeriesAdd())
	mux.HandleFunc("DELETE /api/series/{mId}", a.SeriesRem())

	mux.HandleFunc("PUT /api/series/{mId}/status", a.SeriesUpdateStatus())

	mux.HandleFunc("POST /api/series/episode", a.SeriesEpisodeWatch())
	mux.HandleFunc("PUT /api/series/episode", a.SeriesEpisodeUnWatch())

	// mux.HandleFunc("GET /api/series/{tmdbId}/{episode}", a.SeriesEpisodeGet())

	mux.Handle("GET /", ui.NewSpaHandler("internal/ui/static"))
	return a
}
