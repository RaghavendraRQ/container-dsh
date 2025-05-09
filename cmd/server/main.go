package main

import (
	"container-dsh/pkg/http"
	"log"
)

func main() {
	if err := http.Run(); err != nil {
		log.Fatalln("Failed to start server:", err)
	}
}
