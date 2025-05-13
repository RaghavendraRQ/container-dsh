package http

import (
	"container-dsh/internal/container"
	"encoding/json"
	"log"
	"net/http"
)

func loggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		http.Error(w, "Failed to get container list", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(containers)
	w.WriteHeader(http.StatusOK)
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
