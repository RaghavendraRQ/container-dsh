package container

//go:generate stringer -type=Status -linecomments
type Status uint

const (
	_ Status = iota
	Created
	Running
	Paused
	UnPaused
	Stopped
	Removed
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
	Status string `json:"status"`
}

type ContainerSData struct {
	Containers []Container `json:"containers"`
	Total      int         `json:"total"`
}
