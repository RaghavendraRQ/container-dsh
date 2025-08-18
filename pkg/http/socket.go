package http

import (
	"container-dsh/internal/container"
	"log/slog"
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

	cli := container.GetClient()
	go handleSingleContainer(conn, cli)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	cli := container.GetClient()
	go handleMetrics(conn, cli)
}

func handleMetrics(conn *websocket.Conn, cli *client.Client) {
	defer conn.Close()
	ticker := time.NewTicker(METRICSREFRESHTIME)
	for range ticker.C {
		containerIds, _ := container.GetContainerList(cli)
		if err := conn.WriteJSON(containerIds); err != nil {
			slog.Error("Error writing message", "error", err)
			break
		}
	}
}

func handleSingleContainer(conn *websocket.Conn, cli *client.Client) {
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Info("WebSocket connection closed", "error", err)
				return
			}
			slog.Info("Error reading message", "error", err)
			return
		}
		slog.Info("Received message: %s\n", "", message)

		if messageType == websocket.CloseMessage {
			slog.Info("WebSocket close message received: %s\n", "", message)
			return
		}

		metrics, err := container.GetContainerData(cli, string(message))
		if err != nil {
			slog.Info("Error getting container data:", "", err)
			conn.WriteMessage(websocket.TextMessage, []byte("No data"))
			continue
		}

		ticker := time.NewTicker(METRICSREFRESHTIME)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteJSON(metrics); err != nil {
				slog.Info("Error writing message:", "", err)
				return
			}
			slog.Info("Container data: %s\n", "", metrics.String())
		}
	}
}
