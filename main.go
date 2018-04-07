package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"github.com/takaishi/alpette/server"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name: "server",
			Action: func(c *cli.Context) error {
				server.Start()
				return nil
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
