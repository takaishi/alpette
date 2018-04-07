package client

import (
	"context"
	pb "github.com/takaishi/alpette/protocol"
	"google.golang.org/grpc"
	"log"
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
