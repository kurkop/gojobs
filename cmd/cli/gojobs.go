package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	flags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{Name: "gojobs-api-url"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "gojobs-api-token"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "name"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "image"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "schedule"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "generate_name"}),
		&cli.StringFlag{Name: "config", Value: "gojob.yaml"},
	}

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
					fmt.Println(c.String("name"))
					fmt.Println(c.String("gojobs-api-url"))
					fmt.Println(c.String("gojobs-api-token"))
					fmt.Println(c.String("image"))
					fmt.Println(c.String("schedule"))
					fmt.Println(c.String("generate_name"))
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
		Before:               altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Flags:                flags,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
