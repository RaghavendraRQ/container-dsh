package http

import (
	"container-dsh/internal/config"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Run() error {
	//TODO: I know it looks ugly, need to refactor this...
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("error in configuration: %v", err)
	}
	corsRules := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.CLIENT_URL},
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

	log.Println("Starting HTTP server on port", cfg.PORT)
	return http.ListenAndServe(cfg.PORT, corsHandler)

}

func httpHandler(muxRouter *mux.Router) {
	muxRouter.HandleFunc("/", HomeHandler).Methods(http.MethodGet)
	muxRouter.HandleFunc("/metrics", GetMetric).Methods(http.MethodGet)
	muxRouter.HandleFunc("/metrics/{id}", GetMetricById).Methods(http.MethodGet)
}

func webSocketHandler(webSocketRouter *mux.Router) {
	webSocketRouter.HandleFunc("/", wsHandler)
	webSocketRouter.HandleFunc("/container", wsContainerHandler)
}
