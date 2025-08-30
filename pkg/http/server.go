package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/raghavendrarq/container-dsh/internal/config"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Run() error {
	//TODO: I know it looks ugly, need to refactor this...
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error in configuration: %v", err)
	}
	corsRules := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.Server.ClientUrl},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	muxRouter := mux.NewRouter().StrictSlash(true)
	muxRouter.Use(loggerMiddleWare)

	// HTTP Routes
	httpHandler(muxRouter)

	// WS Routes
	webSocketRouter := muxRouter.PathPrefix("/ws").Subrouter()
	webSocketHandler(webSocketRouter)

	//CORS Handler
	corsHandler := corsRules.Handler(muxRouter)

	log.Println("Starting HTTP server on port", cfg.Server.Port)
	return http.ListenAndServe(cfg.Server.Port, corsHandler)

}

func httpHandler(muxRouter *mux.Router) {
	muxRouter.HandleFunc("/", HomeHandler).Methods(http.MethodGet)
	muxRouter.HandleFunc("/metrics", GetMetric).Methods(http.MethodGet) // TODO: Change the URL since prometheus searches for this scraping
	muxRouter.HandleFunc("/metrics/{id}", GetMetricById).Methods(http.MethodGet)
}

func webSocketHandler(webSocketRouter *mux.Router) {
	webSocketRouter.HandleFunc("/", wsHandler)
	webSocketRouter.HandleFunc("/container", wsContainerHandler)
}
