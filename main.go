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
			},
			Action: func(c *cli.Context) error {
				return server.Start(c)
			},
		},
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				return client.Start()
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
