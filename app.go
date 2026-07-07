package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"gitlab.com/raffleberry/riptvtime/server"
)

func main() {

	// source := metadata.NewTmdbSource()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	mux := http.NewServeMux()

	addr := "127.0.0.1:5667"
	s := server.New(addr, mux)

	// mux.HandleFunc("GET /api/series", apiSeriesAll())
	mux.HandleFunc("GET /api/series/search", apiSeriesSearch())
	mux.HandleFunc("GET /api/series/feed", apiSeriesFeed())
	// mux.HandleFunc("GET /api/series/{tmdbId}", apiSeriesGet())
	// mux.HandleFunc("GET /api/series/{tmdbId}/{episode}", apiSeriesEpisode())
	mux.HandleFunc("POST /api/series/", apiSeriesAdd())
	mux.HandleFunc("POST /api/series/episode/", apiSeriesEpisodeWatch())
	mux.HandleFunc("DELETE /api/series/episode/", apiSeriesEpisodeUnWatch())
	// mux.HandleFunc("DELETE /api/series/{tmdbId}", apiSeriesRemove())
	// mux.HandleFunc("PUT /api/series/{tmdbId}", apiSeriesUpdateStatus())

	mux.Handle("GET /", server.NewSpaHandler("server/ui", "index.html"))

	fmt.Printf("Starting server...\n")

	if err := s.Start(); err != nil {
		panic(err)
	}

	fmt.Printf("Address - http://%s/ ...\n", addr)
	s.WaitSIGINT()

	fmt.Printf("Stopping http://%s/ ...\n", addr)
	s.Stop()

}
