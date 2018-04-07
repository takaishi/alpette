package server

import (
	"github.com/mattn/go-shellwords"
	pb "github.com/takaishi/alpette/protocol"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
)

type taskService struct {
	tasks []*pb.Task
}

func (ts *taskService) Run(c context.Context, p *pb.Task) (*pb.ResponseType, error) {
	for _, task := range ts.tasks {
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
			log.Printf("[INFO] result: %s\n", out)
		}
	}
	return new(pb.ResponseType), nil
}

func Start(c *cli.Context) error {
	log.Println("[DEBUG] server")

	raw, err := ioutil.ReadFile(c.String("conf"))
	if err != nil {
		return err
	}
	var tasks []*pb.Task
	err = yaml.Unmarshal(raw, &tasks)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", ":11111")
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	ts := taskService{tasks: tasks}
	pb.RegisterTaskServiceServer(server, &ts)
	return server.Serve(lis)
}
