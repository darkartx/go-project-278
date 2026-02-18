package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	cli "github.com/urfave/cli/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	cmd := &cli.Command{
		Usage: "Url shortener",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Value:   false,
				Aliases: []string{"D"},
				Usage:   "enable debug mode",
				Sources: cli.EnvVars("DEBUG"),
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "serve api",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "debug",
						Value:   false,
						Aliases: []string{"D"},
						Usage:   "enable debug mode",
						Sources: cli.EnvVars("DEBUG"),
					},
					&cli.StringFlag{
						Name:    "port",
						Value:   "8080",
						Aliases: []string{"p"},
						Usage:   "api port",
						Sources: cli.EnvVars("PORT"),
					},
					&cli.StringFlag{
						Name:    "db",
						Usage:   "database url",
						Sources: cli.EnvVars("DATABASE_URL"),
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					debug := c.Bool("debug")
					port := c.String("port")
					databaseUrl := c.String("db")

					config := NewConfig(debug, databaseUrl, port)

					return Api(config)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
