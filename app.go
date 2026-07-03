package main

import (
	"fmt"
	"net/http"

	"gitlab.com/raffleberry/riptvtime/server"
)

func main() {
	fmt.Println("Hello World")

	mux := http.NewServeMux()

	addr := "127.0.0.1:5667"
	s := server.New(addr, mux)

	mux.Handle("GET /", server.NewSpaHandler("server/ui", "index.html"))

	fmt.Printf("Starting server...\n")

	err := s.Start()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Address - http://%s/ ...\n", addr)
	s.WaitSIGINT()

	fmt.Printf("Stopping http://%s/ ...\n", addr)
	s.Stop()

}
