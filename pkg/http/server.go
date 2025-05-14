package http

import (
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
	httpHandler(muxRouter)

	// WS Routes
	webSocketRouter := muxRouter.PathPrefix("/ws").Subrouter()
	webSocketHandler(webSocketRouter)

	corsHandler := corsRules.Handler(muxRouter)

	log.Println("Starting HTTP server on port", URL)
	return http.ListenAndServe(URL, corsHandler)

}

func httpHandler(muxRouter *mux.Router) {
	muxRouter.HandleFunc("/", HomeHandler).Methods(http.MethodGet)
	muxRouter.HandleFunc("/metrics", GetMetric).Methods(http.MethodGet)
    muxRouter.HandleFunc("/metrics/{id}", GetMetricById).Methods(http.MethodGet)
}

func webSocketHandler(webSocketRouter *mux.Router) {
	webSocketRouter.HandleFunc("", wsHandler)
	webSocketRouter.HandleFunc("/container", wsContainerHandler)
}
