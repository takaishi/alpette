package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"github.com/takaishi/alpette/server"
	"github.com/hashicorp/logutils"
)

func logLevel() string {
	envLevel := os.Getenv("LOG_LEVEL")
	if envLevel == "" {
		return "WARN"
	} else {
		return envLevel
	}
}

func logOutput() (*logutils.LevelFilter){
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(logLevel()),
		Writer: os.Stderr,
	}

	return filter
}

func main() {
	log.SetOutput(logOutput())
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name: "server",
			Action: func(c *cli.Context) error {
				return server.Start()
			},
		},
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
