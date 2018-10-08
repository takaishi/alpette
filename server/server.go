package server

import (
	"fmt"
	"github.com/mattn/go-shellwords"
	pb "github.com/takaishi/alpette/protocol"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/exec"
	"github.com/takaishi/alpette/credentials/stns"
	"github.com/BurntSushi/toml"
)

type taskService struct {
	AuthType string `toml:"authType"`
	Tasks []*pb.Task
}

func (ts *taskService) Run(c context.Context, p *pb.Task) (*pb.ResponseType, error) {
	var resp pb.ResponseType
	for _, task := range ts.Tasks {
		if task.Name == p.Name {
			cmds, err := shellwords.Parse(task.Command)
			if err != nil {
				return nil, err
			}
			log.Printf("[INFO] exec: %s\n", task.Command)

			out, err := exec.Command(cmds[0], cmds[1:]...).Output()
			if err != nil {
				return nil, err
			}
			resp.Body = string(out)
			log.Printf("[INFO] result: %s\n", out)
		}
	}
	return &resp, nil
}

func Start(c *cli.Context) error {
	log.Println("[DEBUG] server")

	var ts taskService
	var server *grpc.Server
	_, err := toml.DecodeFile(c.String("conf"), &ts)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", c.String("port")))
	if err != nil {
		return err
	}

	if ts.AuthType == "STNS" {
		log.Println("[DEBUG] authType is STNS.")
		stnsTC := stns.NewServerCreds(c.String("stns-address"), c.String("stns-port"))
		server = grpc.NewServer(grpc.Creds(stnsTC))
	}

	if ts.AuthType == "insecure" {
		log.Println("[DEBUG] authType is insecure.")
		server = grpc.NewServer()
	}
	pb.RegisterTaskServiceServer(server, &ts)
	return server.Serve(lis)
}
