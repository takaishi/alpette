package main

import (
	"fmt"
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
					Name:  "port",
					Value: "11111",
					Usage: "Sets the gRPC port to listen on.",
				},
				cli.StringFlag{
					Name:  "stns-address",
					Value: "127.0.0.1",
					Usage: "Set the STNS address.",
				},
				cli.StringFlag{
					Name:  "stns-port",
					Value: "1104",
					Usage: "Set the STNS port.",
				},
				cli.StringFlag{
					Name:  "auth",
					Value: "insecure",
					Usage: "Set the authentication method.",
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
					Name:  "server-host",
					Value: "localhost",
					Usage: "alpette server host",
				},
				cli.StringFlag{
					Name:  "server-port",
					Value: "11111",
					Usage: "alpette server port",
				},
				cli.StringFlag{
					Name:  "user, u",
					Value: os.Getenv("USER"),
					Usage: "Set username for authentication.",
				},
				cli.StringFlag{
					Name:  "identity-file, i",
					Value: fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME")),
					Usage: "Set identity (private key) file for authentication.",
				},
				cli.StringFlag{
					Name:  "auth-type",
					Value: "insecure",
					Usage: "Set the authentication method.",
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
