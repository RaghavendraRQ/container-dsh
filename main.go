package main

import (
	"container-dsh/collector"
	"time"
)

func main() {
	collector.Start()
	time.Sleep(5 * time.Second)
}
