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
	mu       sync.Mutex
	buffer   []MetricEntry
	filePath string
}

func (t *TimeSeries) Log(id, metric string, value float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	newMetric := MetricEntry{
		TimeStamp:   time.Now(),
		ContainerId: id,
		Metric:      metric,
		Value:       value,
	}
	t.buffer = append(t.buffer, newMetric)
}

func (t *TimeSeries) Dump(filepath string) {
	t.filePath = filepath
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
