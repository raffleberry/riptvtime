package main

import (
	"fmt"
	"net/http"
	"os"

	tmdb "github.com/cyruzin/golang-tmdb"
	"gitlab.com/raffleberry/riptvtime/server"
)

func main() {
	fmt.Println("Hello World")

	tmdbClient, err := tmdb.Init(os.Getenv("TMDB_API_KEY"))
	if err != nil {
		fmt.Println(err)
	}

	testTmdb := func() {
		movie, err := tmdbClient.GetMovieDetails(297802, nil)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(movie)

		fmt.Println(movie.Title)
	}

	mux := http.NewServeMux()

	addr := "127.0.0.1:5667"
	s := server.New(addr, mux)

	mux.Handle("GET /", server.NewSpaHandler("server/ui", "index.html"))

	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		testTmdb()
		w.Write([]byte("Check Console"))
	})

	fmt.Printf("Starting server...\n")

	err = s.Start()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Address - http://%s/ ...\n", addr)
	s.WaitSIGINT()

	fmt.Printf("Stopping http://%s/ ...\n", addr)
	s.Stop()

}
