package http

import (
	"container-dsh/internal/container"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	METRICSREFRESHTIME = time.Second * 1
)

func wsContainerHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	cli := container.GetClient()
	handleSingleContainer(conn, cli)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	cli := container.GetClient()
	handleMetrics(conn, cli)
}

func handleMetrics(conn *websocket.Conn, cli *client.Client) {
	ticker := time.NewTicker(METRICSREFRESHTIME)
	for range ticker.C {
		containerIds, _ := container.GetContainerList(cli)
		if err := conn.WriteJSON(containerIds); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func handleSingleContainer(conn *websocket.Conn, cli *client.Client) {
	for {
		_, containerId, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message: %s\n", containerId)
		metrics, err := container.GetContainerData(cli, string(containerId))
		if err != nil {
			log.Println("Error getting container data:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("No data"))
			continue
		}
		if err := conn.WriteJSON(metrics); err != nil {
			log.Println("Error writing message:", err)
			break
		}
		log.Printf("Container data: %s\n", metrics.String())
	}
}
