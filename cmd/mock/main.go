package main

import (
	"container-dsh/internal/container"
	"container-dsh/pkg/logger"
	"fmt"
	"sync"
	"time"
)

var (
	timeLogger logger.TimeLogger
	wg         = sync.WaitGroup{}
)

func main() {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		panic(err)
	}
	timeLogger = *logger.NewTimeLogger("time.json")
	go timeLogger.Start()
	for _, containerId := range containers {
		//fmt.Println("Container ID: ", containerId)
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			containerData, err := container.GetContainerData(cli, id)
			if err != nil {
				panic(err)
			}
			sendMetricsToLogger(id, containerData)
			fmt.Println(containerData.String())
		}(containerId)
	}
	wg.Wait()
	timeLogger.Done <- true
	timeLogger.Wg.Wait()

	//time.Sleep(1 * time.Second)

}

func sendMetricsToLogger(id string, metric container.Container) {
	if timeLogger.MetricsChannel != nil {
		timeLogger.MetricsChannel <- logger.MetricEntry{
			TimeStamp:   time.Now(),
			ContainerId: id,
			Metric:      "cpu_usage",
			Value:       metric.CpuUsage,
		}
		timeLogger.MetricsChannel <- logger.MetricEntry{
			TimeStamp:   time.Now(),
			ContainerId: id,
			Metric:      "mem_usage",
			Value:       metric.MemUsage,
		}
		timeLogger.MetricsChannel <- logger.MetricEntry{
			TimeStamp:   time.Now(),
			ContainerId: id,
			Metric:      "net_io",
			Value:       metric.NetIO,
		}
		timeLogger.MetricsChannel <- logger.MetricEntry{
			TimeStamp:   time.Now(),
			ContainerId: id,
			Metric:      "disk_io",
			Value:       metric.DiskIO,
		}
	}
}
