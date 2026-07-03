package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

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

func New(addr string, mux http.Handler) *Server {
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
			log.Println("Error starting server: ", err.Error())
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
			log.Println("Error starting server: ", err.Error())
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
