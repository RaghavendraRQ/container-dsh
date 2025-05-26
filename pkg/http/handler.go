package http

import (
	"container-dsh/internal/container"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func loggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	bytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, "Error in encoding", status)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if _, err := w.Write(bytes); err != nil {
		http.Error(w, "Error in writing the json data", status)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		http.Error(w, "Failed to get container list", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusInternalServerError, containers)
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		http.Error(w, "Can't list containers", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusInternalServerError, containers)
}

func GetMetricById(w http.ResponseWriter, r *http.Request) {
	cli := container.GetClient()
	containerId := mux.Vars(r)["id"]

	stats, err := container.GetContainerData(cli, containerId)

	if err != nil {
		http.Error(w, "Failed to fetch the data", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusInternalServerError, stats)

}
