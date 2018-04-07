package server

import (
	pb "github.com/takaishi/alpette/protocol"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type taskService struct {
	task *pb.Task
}

func (ts *taskService) Run(c context.Context, p *pb.Task) (*pb.ResponseType, error) {
	log.Printf("[DEBUG] task: %#v\n", p)
	return new(pb.ResponseType), nil
}

func Start() error {
	log.Println("[DEBUG] server")

	lis, err := net.Listen("tcp", ":11111")
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterTaskServiceServer(server, new(taskService))
	return server.Serve(lis)
}
