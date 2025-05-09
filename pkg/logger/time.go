package logger

import (
	"encoding/json"
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

type TimeLogger struct {
	sync.WaitGroup
	mu       sync.Mutex
	Buffer   []MetricEntry `json:"timeseries"`
	filePath string
	QuitCh   chan bool
	InputCh  chan MetricEntry
}

func NewTimeLogger(filepath string) *TimeLogger {
	return &TimeLogger{
		filePath: filepath,
		InputCh:  make(chan MetricEntry),
		QuitCh:   make(chan bool),
	}
}

func (t *TimeLogger) Start() {
	t.Add(1)
	defer t.Done()
	for {
		select {
		case <-t.QuitCh:
			log.Print("Done Channel recieived...")
			t.mu.Lock()
			t.Dump()
			t.mu.Unlock()
			return

		case metric, ok := <-t.InputCh:
			if !ok {
				t.Dump()
				return
			}
			t.mu.Lock()
			t.Buffer = append(t.Buffer, metric)
			if len(t.Buffer) >= 100 {
				t.Dump()
				t.Buffer = nil
			}
			t.mu.Unlock()
		}
	}
}

func (t *TimeLogger) Dump() {
	file, err := os.OpenFile(t.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Open file error:", err)
		return
	}
	defer file.Close()

	if len(t.Buffer) == 0 {
		log.Println("Buffer is empty, nothing to write")
		return
	}

	if err := json.NewEncoder(file).Encode(t.Buffer); err != nil {
		log.Println("Error while inserting log:", err)
		return
	}
	//encoder := json.NewEncoder(file)
	//for _, metric := range t.Buffer {
	//if err := encoder.Encode(metric); err != nil {
	//fmt.Printf("Error while inserting log: %s", err)
	//}
	//}
	t.Buffer = nil
}
