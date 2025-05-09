package http

import (
	"container-dsh/internal/container"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	URL string = ":8080"
)

func Run() {
	muxRouter := mux.NewRouter()
	muxRouter.Use(loggerMiddleWare)
	muxRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Container Metrics API"))
	}).Methods("GET")
	muxRouter.HandleFunc("/metrics", GetMetric).Methods("GET")
	log.Println("Starting HTTP server on port", URL)
	http.ListenAndServe(URL, muxRouter)
}

func loggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		panic(err)
	}

	w.Write([]byte("Container IDs: "))
	for _, containerId := range containers {
		w.Write([]byte(containerId + "\n"))
	}
	w.WriteHeader(http.StatusOK)
}
