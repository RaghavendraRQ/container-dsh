package pb

import (
	"fmt"
	"net"

	v1 "github.com/raghavendrarq/container-dsh/api/gen/go/v1"
	"google.golang.org/grpc"
)

const (
	ADDR = ":5001"
	HELP = `
		Addr: :5001
		ContainerServices:
			GetContainerService
	`
)

func Run() error {
	grpcServer := grpc.NewServer()
	listner, err := net.Listen("tcp", ADDR)
	if err != nil {
		return fmt.Errorf("error in creating listner: %v", err)
	}
	v1.RegisterContainerServiceServer(grpcServer, &ContainerService{})

	fmt.Println("GRPC Server is running on Addr: ", ADDR)
	if err := grpcServer.Serve(listner); err != nil {
		return fmt.Errorf("error in serving grpc: %v", err)
	}
	return nil
}
