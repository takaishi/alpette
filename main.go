package main

import (
	"github.com/hashicorp/logutils"
	"github.com/takaishi/alpette/client"
	"github.com/takaishi/alpette/server"
	"github.com/urfave/cli"
	"log"
	"os"
)

func logLevel() string {
	envLevel := os.Getenv("LOG_LEVEL")
	if envLevel == "" {
		return "WARN"
	} else {
		return envLevel
	}
}

func logOutput() *logutils.LevelFilter {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(logLevel()),
		Writer:   os.Stderr,
	}

	return filter
}

func main() {
	log.SetOutput(logOutput())
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name: "server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "conf, c",
					Value: "/path/to/config",
					Usage: "set path to config file.",
				},
				cli.StringFlag{
					Name: "port",
					Value: "11111",
					Usage: "Sets the gRPC port to listen on.",
				},
			},
			Action: func(c *cli.Context) error {
				return server.Start(c)
			},
		},
		{
			Name: "run",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "task",
					Value: "task",
					Usage: " set to exec task name.",
				},
				cli.StringFlag{
					Name: "server-host",
					Value: "localhost",
					Usage: "alpette server host",
				},
				cli.StringFlag{
					Name: "server-port",
					Value: "11111",
					Usage: "alpette server port",
				},
			},
			Action: func(c *cli.Context) error {
				return client.Start(c)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
