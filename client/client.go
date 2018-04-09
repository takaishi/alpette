package client

import (
	"context"
	"fmt"
	pb "github.com/takaishi/alpette/protocol"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"log"
	"github.com/takaishi/alpette/credentials/stns"
)

func Start(c *cli.Context) error {
	log.Println("[DEBUG] client")
	stnsTC := stns.NewClientCreds()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(stnsTC),
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", c.String("server-host"), c.String("server-port")), opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	task := &pb.Task{
		Name: c.String("task"),
	}
	_, err = client.Run(context.Background(), task)
	return err
}
