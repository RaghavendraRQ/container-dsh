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

var (
	corsRules = cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
)

func Run() error {
	muxRouter := mux.NewRouter()
	muxRouter.Use(loggerMiddleWare)

	// HTTP Routes
	muxRouter.HandleFunc("/metrics", GetMetric).Methods("GET")

	// WS Routes
	webSocketRouter := muxRouter.PathPrefix("/ws").Subrouter()
	webSocketHandler(webSocketRouter)

	corsHandler := corsRules.Handler(muxRouter)

	log.Println("Starting HTTP server on port", URL)
	return http.ListenAndServe(URL, corsHandler)

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
