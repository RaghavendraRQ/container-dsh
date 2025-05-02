package main

import (
	"container-dsh/pkg/collector"
	"time"
)

func main() {
	collector.Start(true)

	time.Sleep(2 * time.Second)
}
