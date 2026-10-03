package ui

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/raffleberry/riptvtime/internal/utils"
)

//go:embed static
var ui embed.FS

type spaHandler struct {
	fsys fs.FS
	// baseUrl must begin with / and not end with /
	baseUrl string
}

func replaceAll(fsys fs.FS, filepath string, find []string, replace []string) ([]byte, error) {
	if len(find) != len(replace) {
		return nil, errors.New("find and replace must be same length")
	}
	file, err := fsys.Open(filepath)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	contentStr := string(content)

	for i := range find {
		contentStr = strings.ReplaceAll(contentStr, find[i], replace[i])
	}

	return []byte(contentStr), nil

}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	writeHttpErr := func(err error) {
		http.Error(w, fmt.Errorf("%v, path: %v", err, r.URL.Path).Error(), http.StatusInternalServerError)
	}

	serveIndex := func() {
		slog.Debug("replacing prefix", "baseUrl", h.baseUrl)
		content, err := replaceAll(h.fsys, "index.html",
			[]string{
				`href="/`,
				`src="/`,
			},
			[]string{
				fmt.Sprintf(`href="%s/`, h.baseUrl),
				fmt.Sprintf(`src="%s/`, h.baseUrl),
			},
		)
		if err != nil {
			writeHttpErr(err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, err = w.Write(content)
		if err != nil {
			writeHttpErr(err)
			return
		}
	}

	p := path.Clean(r.URL.Path)

	slog.Debug("UI request", "requestPath", p)

	p = strings.TrimPrefix(p, "/")

	slog.Debug("UI request", "Trimmed requestPath", p)

	if p == "" || p == "index.html" {
		serveIndex()
		return
	}

	info, err := fs.Stat(h.fsys, p)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			serveIndex()
			return
		}

		writeHttpErr(err)
		return
	}

	if info.IsDir() {
		serveIndex()
		return
	}

	if p == "app/base.js" {
		baseJs, err := h.fsys.Open(p)
		if err != nil {
			writeHttpErr(err)
			return
		}

		content, err := io.ReadAll(baseJs)
		if err != nil {
			writeHttpErr(err)
			return
		}

		content, err = replaceAll(h.fsys, p,
			[]string{
				`""`,
			},
			[]string{
				`"` + h.baseUrl + `"`,
			},
		)
		if err != nil {
			writeHttpErr(err)
			return
		}

		w.Header().Set("Content-Type", "application/javascript")
		_, err = w.Write(content)
		if err != nil {
			writeHttpErr(err)
			return
		}

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

// path: abs Path from pwd to static folder from root of this code repo
//
// used during dev
//
// baseUrl must begin with / and not end with / or must be empty
func NewSpaHandler(path string, baseUrl string) http.Handler {

	if len(baseUrl) > 0 {
		if !strings.HasPrefix(baseUrl, "/") {
			panic("baseUrl must begin with /")
		}
		if strings.HasSuffix(baseUrl, "/") {
			panic("baseUrl must not end with /")
		}
	}

	slog.Debug("spa handler", "baseUrl", baseUrl)

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
	return &spaHandler{fsys, baseUrl}
}
