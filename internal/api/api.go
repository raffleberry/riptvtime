package api

import (
	"log/slog"
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

	mux.HandleFunc("POST /api/series/episode", a.SeriesEpisodeWatch())
	mux.HandleFunc("PUT /api/series/episode", a.SeriesEpisodeUnWatch())

	mux.HandleFunc("POST /api/import/upload", a.SeriesImportUpload())
	mux.HandleFunc("POST /api/import/process", func(w http.ResponseWriter, r *http.Request) {
		// TODO
		go func() {
			if state.Import.GetUploadActive() {
				return
			}
			state.Import.SetUploadActive(true)
			defer state.Import.SetUploadActive(false)
			state.Import.Reset()

			err := a.tv.IptImportTvTimeData("")
			if err != nil {
				slog.Error("Error while importing tv time series", "error", err)
				state.Import.SetUploadError(err)
			}
		}()
	})
	mux.HandleFunc("GET /api/import/list", a.SeriesImportDataUnresolved())
	mux.HandleFunc("PUT /api/import/match", a.SeriesImportMatchAndRemove())

	mux.HandleFunc("GET /api/state", a.GetState())
	// mux.HandleFunc("GET /api/series/{tmdbId}/{episode}", a.SeriesEpisodeGet())

	mux.Handle("GET /", ui.NewSpaHandler("internal/ui/static"))
	return a
}
