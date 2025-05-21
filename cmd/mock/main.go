package logger

import (
	"container-dsh/internal/container"
	"container-dsh/pkg/aggr"
	"container-dsh/pkg/logger"
	"fmt"
	"sync"
	"time"
)

var (
	timeLogger *logger.TimeLogger
	manager    *aggr.AggregatorManager
	wg         = sync.WaitGroup{}
)

func Run() {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		panic(err) // TODO: Want to propagate the error
	}
	manager = aggr.NewAggregatorManager(2, 3)
	go manager.Run()
	timeLogger = logger.NewTimeLogger("time.json")
	go timeLogger.Start()
	for _, containerId := range containers {
		//fmt.Println("Container ID: ", containerId)
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			containerData, err := container.GetContainerData(cli, id)
			if err != nil {
				panic(err) // TODO: Want to propagate the error
			}
			sendMetricsToLogger(id, containerData)
			manager.Input <- containerData
			fmt.Println(containerData.String())
		}(containerId)
	}
	wg.Wait()
	timeLogger.QuitCh <- true
	timeLogger.Wait()
	time.Sleep(5 * time.Second)
	manager.Stop()

	//time.Sleep(1 * time.Second)

}

func sendMetricsToLogger(id string, metric container.Container) {
	if timeLogger.InputCh != nil {
		timeLogger.InputCh <- logger.MetricEntry{
			TimeStamp:   time.Now(),
			ContainerId: id,
			Metric:      "cpu_usage",
			Value:       metric.CpuUsage,
		}
		timeLogger.InputCh <- logger.MetricEntry{
			TimeStamp:   time.Now(),
			ContainerId: id,
			Metric:      "mem_usage",
			Value:       metric.MemUsage,
		}
		timeLogger.InputCh <- logger.MetricEntry{
			TimeStamp:   time.Now(),
			ContainerId: id,
			Metric:      "net_io",
			Value:       metric.NetIO,
		}
		timeLogger.InputCh <- logger.MetricEntry{
			TimeStamp:   time.Now(),
			ContainerId: id,
			Metric:      "disk_io",
			Value:       metric.DiskIO,
		}
	}
}
