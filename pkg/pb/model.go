package pb

import (
	"github.com/gogo/protobuf/proto"
	v1 "github.com/raghavendrarq/container-dsh/api/gen/go/v1"
	"github.com/raghavendrarq/container-dsh/internal/container"
)

type ContainerStat struct {
	CpuUsage float64
	MemUsage float64
	NetIO    float64
	DiskIO   float64
}

func (c *ContainerStat) ToProto() *v1.ContainerStat {
	return &v1.ContainerStat{
		CpuUsage: proto.Float64(c.CpuUsage),
		MemUsage: proto.Float64(c.MemUsage),
		Netio:    proto.Float64(c.NetIO),
		Diskio:   proto.Float64(c.DiskIO),
	}
}

func NewContainerStats(cont container.Container) ContainerStat {
	return ContainerStat{
		CpuUsage: cont.CpuUsage,
		MemUsage: cont.MemUsage,
		NetIO:    cont.NetIO,
		DiskIO:   cont.DiskIO,
	}
}
