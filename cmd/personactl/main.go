package main

import (
	"fmt"
	cli "github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "personactl",
		Version:     "0.1.0",
		Description: "a cli application to moderate persona instance.",
		Commands: []*cli.Command{
			{
				Name:    "migrate",
				Aliases: []string{"db"},
				Usage:   "Execute migration for database.",
				Subcommands: []*cli.Command{
					{
						Name:    "down",
						Aliases: []string{"back", "b"},
						Usage:   "Roll back migrations",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:  "times, t",
								Value: 1,
							},
						},
						Action: func(context *cli.Context) error {
							fmt.Printf("roll backed %d migrations.", context.Int("times"))
							return nil
						},
					},
					{
						Name:    "up",
						Aliases: []string{"u"},
						Usage:   "",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:  "times, t",
								Value: 1,
							},
						},
						Action: func(context *cli.Context) error {
							fmt.Printf("applied %d migrations.", context.Int("times"))
							return nil
						},
					},
				},
			},
		},
	}
	_ = app.Run(os.Args)
}
