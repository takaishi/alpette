package client

import (
	"log"
	"google.golang.org/grpc"
	pb "github.com/takaishi/alpette/protocol"
	"context"
)

func Start() error {
	log.Println("[DEBUG] client")
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("127.0.0.1:11111", opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	task := &pb.Task{
		Name: "foo",
	}
	_, err = client.Run(context.Background(), task)
	return err
}