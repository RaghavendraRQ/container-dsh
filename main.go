package main

import (
	"container-dsh/collector"
	"log"
	"time"
)

func main() {
	collector.Start(true)
	log.Println("After collecting")
	time.Sleep(5 * time.Second)
}
