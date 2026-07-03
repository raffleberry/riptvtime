package server

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"syscall"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func WithCtx(handler func(*Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{W: w, R: r}
		err := handler(ctx)
		if err != nil && !errors.Is(err, syscall.EPIPE) {
			log.Printf("%s - Error: %v\n", r.RequestURI, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (c *Context) JSON(status int, data interface{}) error {
	c.W.Header().Set("Content-Type", "application/json")
	c.W.WriteHeader(status)
	return json.NewEncoder(c.W).Encode(data)
}

func (c *Context) File(p string) error {
	http.ServeFile(c.W, c.R, p)
	return nil
}

func (c *Context) FileFromFS(filePath string, fsys fs.FS) error {
	http.ServeFileFS(c.W, c.R, fsys, filePath)
	return nil
}

func (c *Context) Error(status int, message string) error {
	http.Error(c.W, message, status)
	return nil
}
