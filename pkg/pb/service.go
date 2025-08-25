package pb

import (
	"context"
	"log"

	v1 "github.com/raghavendrarq/container-dsh/api/gen/go/v1"
	"github.com/raghavendrarq/container-dsh/internal/container"
)

type ContainerService struct {
	v1.UnimplementedContainerServiceServer
}

func (c *ContainerService) GetContainerStats(ctx context.Context, req *v1.ContainerStatRequest) (*v1.Container, error) {
	cli := container.GetClient()
	containerData, err := container.GetContainerData(cli, "3cd562e2b938")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	Container := NewContainer(containerData)
	ContainerPb := Container.ToProto()
	return ContainerPb, nil
}
