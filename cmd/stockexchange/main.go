package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":9292", "Listen and serve HTTP on host:port")
}

func main() {
	flag.Parse()

	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatalf("The provided %s addr is not correct format (ex. host:port)", addr)
	}

	router := mux.NewRouter()
	router.Handle("/", http.NotFoundHandler())

	server := negroni.New()
	server.Use(negroni.NewRecovery())
	server.UseHandler(router)

	log.Printf("StackExchange started. HTTP listen and serve on %s", addr)
	http.ListenAndServe(addr, router)
}
