package server

import (
	"container-dsh/pkg/http"
	"log"
)

func Run() {
	if err := http.Run(); err != nil {
		log.Fatalln("Failed to start server:", err)
	}
}
