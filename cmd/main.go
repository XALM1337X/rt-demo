package main

import (
	"log"

	"github.com/XALM1337X/rt-demo/internal/server"
	_ "github.com/lib/pq"
)

func main() {
	srv := server.NewHTTPServer(":6543")
	log.Fatal(srv.ListenAndServe())
}
