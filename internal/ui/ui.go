package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"gitlab.com/raffleberry/riptvtime/internal/utils"
)

//go:embed static
var ui embed.FS

type spaHandler struct {
	serveDir  string
	indexFile string
	fsys      fs.FS
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := filepath.Join(h.serveDir, filepath.Clean(r.URL.Path))

	if info, err := fs.Stat(h.fsys, p); err != nil {
		http.ServeFile(w, r, filepath.Join(h.serveDir, h.indexFile))
		return
	} else if info.IsDir() {
		http.ServeFile(w, r, filepath.Join(h.serveDir, h.indexFile))
		return
	}

	http.ServeFile(w, r, p)
}

func NewSpaHandler() http.Handler {
	publicDir := "static"
	indexFile := "index.html"
	var fsys fs.FS
	var err error
	if utils.IsGoRun() {
		fsys = os.DirFS(".")
		fmt.Println("UI - using live mode")
	} else {
		fmt.Println("UI - using embed mode")
		fsys, err = fs.Sub(ui, "ui")
		if err != nil {
			panic(err)
		}
	}
	return &spaHandler{publicDir, indexFile, fsys}
}
