package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:            "gojob",
		Usage:           "a tool to publish and run jobs on Kubernetes",
		UsageText:       "gojob [global options] command [command options] [arguments...]",
		HideHelpCommand: true,
		Compiled:        time.Now(),
		Commands: []*cli.Command{
			{
				Name:    "publish",
				Aliases: []string{"p"},
				Usage:   "Build a new docker image and upload it to registry",
				Action: func(c *cli.Context) error {
					fmt.Println("Publishing... ", c.Args().First())
					return nil
				},
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Create a new job or cronjob on Kubernetes",
				Action: func(c *cli.Context) error {
					fmt.Println("Running", c.Args().First())
					return nil
				},
			},
		},
		EnableBashCompletion: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
