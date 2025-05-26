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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		http.Error(w, "Failed to get container list", http.StatusInternalServerError)
		return
	}

	bytes, err := json.MarshalIndent(containers, "", " ")
	if err != nil {
		http.Error(w, "Error in encoding", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if _, err := w.Write(bytes); err != nil {
		http.Error(w, "Error in writing the json data", http.StatusInternalServerError)
		return
	}
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		http.Error(w, "Can't list containers", http.StatusInternalServerError)
		return
	}
	data, err := json.MarshalIndent(containers, "", " ")
	if err != nil {
		http.Error(w, "Error in serialising the data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if _, err := w.Write(data); err != nil {
		http.Error(w, "Error in writing the json data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetMetricById(w http.ResponseWriter, r *http.Request) {
	cli := container.GetClient()
	containerId := mux.Vars(r)["id"]

	stats, err := container.GetContainerData(cli, containerId)

	if err != nil {
		http.Error(w, "Failed to fetch the data", http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(stats)

	if err != nil {
		http.Error(w, "Error in encoding", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	if _, err := w.Write(bytes); err != nil {
		http.Error(w, "Error in writing the json data", http.StatusInternalServerError)
		return
	}
}
