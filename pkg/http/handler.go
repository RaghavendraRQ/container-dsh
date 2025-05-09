package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	URL string = ":8080"
)

func Run() {
	muxRouter := mux.NewRouter()
	log.Println("Starting HTTP server on port", URL)
	http.ListenAndServe(URL, muxRouter)
}
