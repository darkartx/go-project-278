package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	config, err := getConfigFromEnv()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	err = Api(&config)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func getConfigFromEnv() (Config, error) {
	result := Config{
		Debug:       false,
		DatabaseUrl: "",
		Bind:        "127.0.0.1:8080",
	}

	if debugEnv, exists := os.LookupEnv("DEBUG"); exists {
		debug, err := strconv.ParseBool(debugEnv)
		if err != nil {
			return Config{}, err
		}

		result.Debug = debug
	}

	if databaseUrl, exists := os.LookupEnv("DATABASE_URL"); exists {
		result.DatabaseUrl = databaseUrl
	}

	if bind, exists := os.LookupEnv("BIND"); exists {
		result.Bind = bind
	}

	return result, nil
}

// func main() {
// 	if err := godotenv.Load(); err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 	}

// 	cmd := &cli.Command{
// 		Usage: "Url shortener",
// 		Flags: []cli.Flag{
// 			&cli.BoolFlag{
// 				Name:    "debug",
// 				Value:   false,
// 				Aliases: []string{"D"},
// 				Usage:   "enable debug mode",
// 				Sources: cli.EnvVars("DEBUG"),
// 			},
// 		},
// 		Commands: []*cli.Command{
// 			{
// 				Name:    "serve",
// 				Aliases: []string{"s"},
// 				Usage:   "serve api",
// 				Flags: []cli.Flag{
// 					&cli.BoolFlag{
// 						Name:    "debug",
// 						Value:   false,
// 						Aliases: []string{"D"},
// 						Usage:   "enable debug mode",
// 						Sources: cli.EnvVars("DEBUG"),
// 					},
// 					&cli.StringFlag{
// 						Name:    "port",
// 						Value:   "8080",
// 						Aliases: []string{"p"},
// 						Usage:   "api port",
// 						Sources: cli.EnvVars("PORT"),
// 					},
// 					&cli.StringFlag{
// 						Name:    "db",
// 						Usage:   "database url",
// 						Sources: cli.EnvVars("DATABASE_URL"),
// 					},
// 				},
// 				Action: func(ctx context.Context, c *cli.Command) error {
// 					debug := c.Bool("debug")
// 					port := c.String("port")
// 					databaseUrl := c.String("db")

// 					config := NewConfig(debug, databaseUrl, port)

// 					return Api(config)
// 				},
// 			},
// 		},
// 	}

// 	if err := cmd.Run(context.Background(), os.Args); err != nil {
// 		fmt.Fprintln(os.Stderr, "Error:", err)
// 		os.Exit(1)
// 	}
// }
