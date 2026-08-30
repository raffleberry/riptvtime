package api

import (
	"encoding/json"
	"net/http"

	"github.com/raffleberry/riptvtime/internal/config"
	"github.com/raffleberry/riptvtime/internal/db"
	"github.com/raffleberry/riptvtime/internal/meta"
	"github.com/raffleberry/riptvtime/internal/services"
	"github.com/raffleberry/riptvtime/internal/services/state"
	"github.com/raffleberry/riptvtime/internal/ui"
)

// if response is not 200, error should be sent in plain text as response
// with the appropriate error code as status
type Api struct {
	Router *http.ServeMux
	db     db.Db
	meta   meta.Meta
	tv     *services.SeriesService
	cfg    *config.Config
}

func NewApi(db db.Db, meta meta.Meta, tv *services.SeriesService, cfg *config.Config) *Api {
	mux := http.NewServeMux()

	a := &Api{
		Router: mux,
		db:     db,
		meta:   meta,
		tv:     tv,
		cfg:    cfg,
	}

	mux.HandleFunc("GET /api/series", a.SeriesAll())
	mux.HandleFunc("GET /api/series/search", a.SeriesSearch())
	mux.HandleFunc("GET /api/series/feed", a.SeriesFeed())

	mux.HandleFunc("GET /api/series/{mId}", a.SeriesGet())
	mux.HandleFunc("GET /api/series/{mId}/upnext", a.SeriesUpNext())
	mux.HandleFunc("POST /api/series", a.SeriesAdd())
	mux.HandleFunc("DELETE /api/series/{mId}", a.SeriesRem())

	mux.HandleFunc("PUT /api/series/{mId}/status", a.SeriesUpdateStatus())
	mux.HandleFunc("GET /api/series/{mId}/poster", a.SeriesPoster())

	mux.HandleFunc("POST /api/series/episode", a.SeriesEpisodeWatch())
	mux.HandleFunc("PUT /api/series/episode", a.SeriesEpisodeUnWatch())
	mux.HandleFunc("GET /api/series/stats", a.SeriesStats())
	mux.HandleFunc("GET /api/series/stats/my", a.SeriesStatsMy())

	mux.HandleFunc("GET /api/series/upcoming", a.SeriesUpcoming())

	mux.HandleFunc("POST /api/import/upload", a.SeriesImportUpload())
	mux.HandleFunc("GET /api/import/unresolved", a.SeriesImportUnresolved())
	mux.HandleFunc("PUT /api/import/resolve", a.SeriesImportResolve())

	mux.HandleFunc("GET /api/state", a.GetState())
	// mux.HandleFunc("POST /api/state", a.SetState()) // dev
	mux.HandleFunc("DELETE /api/state", a.ResetState())
	// mux.HandleFunc("GET /api/series/{tmdbId}/{episode}", a.SeriesEpisodeGet())

	mux.Handle("GET /", ui.NewSpaHandler("internal/ui/static"))
	return a
}

func (a *Api) SetState() http.HandlerFunc {
	return WithCtx(func(ctx *Context) error {
		mapState := make(map[string]any)
		d := json.NewDecoder(ctx.R.Body)
		if err := d.Decode(&mapState); err != nil {
			return err
		}

		state.Import.Set(mapState)
		return ctx.JSON(http.StatusOK, struct{}{})
	})
}
