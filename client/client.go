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
	stnsTC := stns.NewClientCreds(c.String("user"))
	var opts []grpc.DialOption

	authType := c.String("auth-type")
	if authType == "insecure" {
		opts = []grpc.DialOption{
			grpc.WithInsecure(),
		}
	}

	if authType == "stns" {
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(stnsTC),
		}
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", c.String("host"), c.String("port")), opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	task := &pb.Task{
		Name: c.String("task"),
	}
	resp, err := client.Run(context.Background(), task)
	if err != nil {
		return err
	}
	fmt.Printf("%s", resp.Body)

	return nil
}
