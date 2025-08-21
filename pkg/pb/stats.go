package pb

import (
	"log"

	"github.com/raghavendrarq/container-dsh/internal/container"
)

func Run() error {
	cli := container.GetClient()
	containerData, err := container.GetContainerData(cli, "3cd562e2b938")
	if err != nil {
		log.Println(err)
		return err
	}
	ContainerStat := NewContainerStats(containerData)
	ContainerStatPb := ContainerStat.ToProto()
	log.Println(ContainerStatPb.GetCpuUsage())
	return nil
}
