package server

import (
	"log"

	"github.com/raghavendrarq/container-dsh/pkg/http"
)

func Run() {
	if err := http.Run(); err != nil {
		log.Fatalln("Failed to start server:", err)
	}
}
