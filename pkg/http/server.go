package http

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	URL        string
	CLIENT_URL string
)

func init() {
	var ok bool
	URL, ok = os.LookupEnv("PORT")

	if !ok {
		log.Fatalf("Error Uninitialised URL:PORT")
	}

	CLIENT_URL, ok = os.LookupEnv("CLIENT_URL")
	if !ok {
		log.Fatalf("Error Uninitialised CLIENT_URL:PORT")
	}
}

var (
	corsRules = cors.New(cors.Options{
		AllowedOrigins:   []string{CLIENT_URL},
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

	//CORS Handler
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
