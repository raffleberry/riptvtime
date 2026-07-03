package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gitlab.com/raffleberry/riptvtime/metadata"
	"gitlab.com/raffleberry/riptvtime/server"
)

func main() {
	fmt.Println("Hello World")

	source := metadata.NewTmdbSource()

	mux := http.NewServeMux()

	addr := "127.0.0.1:5667"
	s := server.New(addr, mux)

	mux.Handle("GET /", server.NewSpaHandler("server/ui", "index.html"))

	mux.HandleFunc("GET /search", server.WithCtx(func(c *server.Context) error {
		query := c.R.URL.Query().Get("query")
		pageStr := c.R.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Printf("Invalid Request, `page` should be an int, but received '%v'\n", page)
			page = 1
		}
		res, err := source.SearchShows(query, page)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, res)
	}))

	fmt.Printf("Starting server...\n")

	if err := s.Start(); err != nil {
		panic(err)
	}

	fmt.Printf("Address - http://%s/ ...\n", addr)
	s.WaitSIGINT()

	fmt.Printf("Stopping http://%s/ ...\n", addr)
	s.Stop()

}
