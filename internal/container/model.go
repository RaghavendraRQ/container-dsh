package container

import "fmt"

//go:generate stringer -type=Status
type Status uint

// created", "running", "paused", "restarting", "removing", "exited", or "dead
var (
	statusMap = map[string]Status{
		"created":    Created,
		"running":    Running,
		"paused":     Paused,
		"restarting": Restarting,
		"removing":   Removing,
		"exited":     Exited,
		"dead":       Dead,
	}
)

const (
	_ Status = iota
	Created
	Running
	Paused
	Restarting
	Removing
	Exited
	Dead
)

type Stats struct {
	CpuUsage float64 `json:"cpu_usage"`
	MemUsage float64 `json:"mem_usage"`
	NetIO    float64 `json:"net_io"`
	DiskIO   float64 `json:"disk_io"`
}

type Container struct {
	Stats
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"`
}

type ContainerSData struct {
	Containers []Container `json:"containers"`
	Total      int         `json:"total"`
}

func (s *Stats) String() string {
	return fmt.Sprint("CPU Usage: ", s.CpuUsage, " Mem Usage: ", s.MemUsage, " Net IO: ", s.NetIO, " Disk IO: ", s.DiskIO)
}

func (c *Container) String() string {
	return fmt.Sprint("ID: ", c.ID[:10], " Name: ", c.Name, " Status:", c.Status.String(), " Stats: [", c.Stats.String(), "]")
}
