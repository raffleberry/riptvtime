package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/raffleberry/riptvtime/internal/utils"
)

//go:embed static
var ui embed.FS

type spaHandler struct {
	fsys fs.FS
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := filepath.Join(filepath.Clean(r.URL.Path))

	p = strings.TrimPrefix(p, "/")

	if info, err := fs.Stat(h.fsys, p); err != nil || info.IsDir() {
		http.ServeFileFS(w, r, h.fsys, filepath.Join("index.html"))
		return
	}

	http.ServeFileFS(w, r, h.fsys, p)
}

func debugFsys(fsys fs.FS) {
	_ = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			fmt.Println("::::Debug Fsys:::: - " + path)
		}
		return nil
	})

}

// path: Path to static folder from root of this code repo
func NewSpaHandler(path string) http.Handler {
	var fsys fs.FS
	var err error
	if utils.IsGoRun() {
		fsys = os.DirFS(path)
		slog.Info("UI - using live mode")
	} else {
		slog.Info("UI - using embed mode")
		fsys, err = fs.Sub(ui, "static")
		if err != nil {
			panic(err)
		}
	}

	// debugFsys(fsys)

	if err != nil {
		log.Fatal(err)
	}
	return &spaHandler{fsys}
}
