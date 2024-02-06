package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cbebe/startup/pkg"
)

func startServer() {
	port := 9000
	r := startup.Handler()
	s := http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port), Handler: r}
	log.Printf("listening on %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func main() {
	startServer()
}
