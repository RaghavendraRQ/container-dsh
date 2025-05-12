package http

import (
	"container-dsh/internal/container"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	URL string = ":8080"
)

func Run() error {
	muxRouter := mux.NewRouter()
	muxRouter.Use(loggerMiddleWare)

	muxRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Container Metrics API"))
	}).Methods("GET")
	log.Println("Starting HTTP server on port", URL)

	webSocketRouter := muxRouter.PathPrefix("/ws").Subrouter()
	webSocketHandler(webSocketRouter)

	muxRouter.HandleFunc("/metrics", GetMetric).Methods("GET")

    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, 
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
 
    })
    corsHandler := c.Handler(muxRouter)

	return http.ListenAndServe(URL, corsHandler)

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
