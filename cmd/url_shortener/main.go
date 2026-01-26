package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/urfave/cli/v3"

	"code"
)

func main() {
	cmd := &cli.Command{
		Usage:     "Compares two configuration files and shows a difference.",
		Flags:     []cli.Flag{},
		Arguments: []cli.Argument{},
		Commands:  []*cli.Command{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return code.Api()
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
