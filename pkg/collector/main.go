package collector

import (
	"container-dsh/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Metrics struct {
	CPUUsage float64 `json:"cpu_usage"`
	MemUsage float64 `json:"mem_usage"`
	NetIO    float64 `json:"net_io"`
	DiskIO   float64 `json:"disk_io"`
}

const (
	TIMEFILE   string        = "time.json"
	TIMELOGEND time.Duration = 1 * time.Second
)

var (
	ctx        context.Context
	statsPool  sync.Pool
	timeLogger logger.TimeLogger
	wg         sync.WaitGroup
)

func init() {
	// Initialize the collector
	ctx = context.Background()
	statsPool = sync.Pool{
		New: func() any {
			return &container.StatsResponse{}
		},
	}
}

func Start(haveTimeLogger bool) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error in creating client: %v", err)
	}
	defer cli.Close()
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("error in getting containers: %v", err)
	}
	if haveTimeLogger {
		wg.Add(1)
		go startTimeLogger(&wg)
	}
	GetContainerData(cli, containers)
	log.Printf("Before wait")
	wg.Wait()
	log.Printf("After wait")
	return nil
}

func startTimeLogger(wg *sync.WaitGroup) {
	timeLogger = *logger.NewTimeLogger(TIMEFILE)
	go timeLogger.Start()
	go func() {
		defer wg.Done()
		startTime := time.Now()
		for {
			if time.Since(startTime) >= TIMELOGEND {
				timeLogger.Done <- true
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func GetContainerData(cli *client.Client, containers []container.Summary) {
	for _, cont := range containers {

		go func(id string) {
			stats, err := cli.ContainerStatsOneShot(ctx, id)
			if err != nil {
				fmt.Println("Error in ContainerStatsOneShot: ", err)
				return
			}
			defer stats.Body.Close()
			statsData := statsPool.Get().(*container.StatsResponse)
			if err := json.NewDecoder(stats.Body).Decode(statsData); err != nil {
				fmt.Println("Decode err: ", err)
				return
			}

			metrics := NewMetrics(statsData)
			sendMetricsToLogger(id, metrics)
			PrintPretty(id, metrics)

			statsPool.Put(statsData)
		}(cont.ID)
	}
}

func sendMetricsToLogger(id string, metric Metrics) {
	if timeLogger.MetricsChannel != nil {
		timeLogger.MetricsChannel <- logger.MetricEntry{
			TimeStamp:   time.Now(),
			ContainerId: id,
			Metric:      "cpu_usage",
			Value:       metric.CPUUsage,
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

func NewMetrics(statsData *container.StatsResponse) Metrics {
	return Metrics{
		CPUUsage: calculateCPUPercent(statsData),
		MemUsage: calculateMemUsage(statsData),
		NetIO:    calculateNetIO(statsData),
		DiskIO:   calculateDiskIO(statsData),
	}

}

func calculateCPUPercent(stat *container.StatsResponse) float64 {
	cpuDelta := float64(stat.CPUStats.CPUUsage.TotalUsage - stat.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stat.CPUStats.SystemUsage - stat.PreCPUStats.SystemUsage)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		return (cpuDelta / systemDelta) * float64(len(stat.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return 0.0
}

func calculateMemUsage(stat *container.StatsResponse) float64 {
	return float64(stat.MemoryStats.Usage) / (1024 * 1024)
}

func calculateNetIO(stat *container.StatsResponse) float64 {
	if len(stat.Networks) > 0 {
		for _, net := range stat.Networks {
			return float64(net.RxBytes+net.TxBytes) / (1024 * 1024)
		}
	}
	return 0.0
}
func calculateDiskIO(stat *container.StatsResponse) float64 {
	if len(stat.BlkioStats.IoServiceBytesRecursive) > 0 {
		return float64(stat.BlkioStats.IoServiceBytesRecursive[0].Value) / (1024 * 1024)
	}
	return 0.0
}
