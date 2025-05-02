package aggregator

import (
	"container-dsh/logger"
)

type CombinedAggregator struct {
	MetricsChannel chan struct {
		ContainerId string
		Metrics     logger.ContainerMetrics
	}
}

func BlackBox() {
	combinedAggregator := CombinedAggregator{
		MetricsChannel: make(chan struct {
			ContainerId string
			Metrics     logger.ContainerMetrics
		}),
	}
	aggregator := NewAggregator(10)
	aggregator.Start()
	if aggregator == nil {
		panic("Aggregator is nil")
	}

	for {
		select {
		case metric := <-combinedAggregator.MetricsChannel:
			if metric.ContainerId == "0de5244ede" {
				aggregator.MetricsChannel <- metric.Metrics
			}
		}
	}

}
