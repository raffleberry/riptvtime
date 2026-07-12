package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func WithCtx(handler func(*Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request", "uri", r.RequestURI, "method", r.Method)
		ctx := &Context{W: w, R: r}
		err := handler(ctx)
		if err != nil && !errors.Is(err, syscall.EPIPE) {
			slog.Error("Error while handling request", "uri", r.RequestURI, "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (c *Context) JSON(status int, data any) error {
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

// when Stop() is called, the server will wait for this duration before shutting down ...Ungracefully
const StopTimeout = 5 * time.Second

type Server struct {
	server http.Server

	isHttps bool

	ready chan bool
	done  chan bool

	addr string
	mux  http.Handler
}

func NewServer(addr string, mux http.Handler) *Server {
	s := &Server{
		addr: addr,
		mux:  mux,
	}

	s.ready = make(chan bool, 1)
	s.done = make(chan bool, 1)
	s.server = http.Server{Handler: mux}
	return s
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp4", s.addr)
	if err != nil {
		return err
	}

	go func() {
		defer close(s.ready)
		defer close(s.done)

		s.ready <- true

		err = s.server.Serve(listener)

		if err != nil && err != http.ErrServerClosed {
			slog.Error("Error starting server", "err", err)
		}
	}()

	<-s.ready

	return nil
}

func (s *Server) StartHttps(tlsConfig *tls.Config) error {
	listener, err := net.Listen("tcp4", s.addr)
	if err != nil {
		return err
	}

	s.isHttps = true
	s.server.TLSConfig = tlsConfig

	go func() {
		defer close(s.ready)
		defer close(s.done)

		s.ready <- true

		err = s.server.ServeTLS(listener, "", "")

		if err != nil && err != http.ErrServerClosed {
			slog.Error("Error starting server", "err", err)
		}
	}()

	<-s.ready

	return nil
}

func (s *Server) StopWithCtx(ctx context.Context) error {

	s.server.Shutdown(ctx)

	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout")
	}

}

func (s *Server) Addr() string {
	proto := "http"
	if s.isHttps {
		proto = "https"
	}
	return fmt.Sprintf("%s://%s/", proto, s.addr)
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(StopTimeout))
	defer cancel()

	s.server.Shutdown(ctx)

	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout")
	}
}

func (s *Server) WaitSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
