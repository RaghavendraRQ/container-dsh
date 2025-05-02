package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type MetricEntry struct {
	TimeStamp   time.Time `json:"timestamp"`
	ContainerId string    `json:"container_id"`
	Metric      string    `json:"metric"`
	Value       float64   `json:"value"`
}

type TimeSeries struct {
	mu             sync.Mutex
	buffer         []MetricEntry
	filePath       string
	Done           chan bool
	MetricsChannel chan MetricEntry
}

func (t *TimeSeries) Start(filepath string) {
	t.filePath = filepath
	for {
		select {
		case <-t.Done:
			t.Dump()
			return

		case metric := <-t.MetricsChannel:
			t.mu.Lock()
			t.buffer = append(t.buffer, metric)
			if len(t.buffer) >= 100 {
				t.Dump()
				t.buffer = nil
			}
			t.mu.Unlock()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (t *TimeSeries) Dump() {
	file, err := os.OpenFile(t.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Open file error:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, metric := range t.buffer {
		if err := encoder.Encode(metric); err != nil {
			fmt.Printf("Error while inserting log: %s", err)
		}
	}
	t.buffer = nil
}
