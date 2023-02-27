package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/miy4/archii"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}

	fs := flag.NewFlagSet("archii", flag.ExitOnError)
	dir := fs.String("d", ".", "output directory")

	root := &ffcli.Command{
		ShortUsage: "archii [-d dir] <URL>",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires an URL to read")
			}
			return archii.RunApp(args[0], *dir)
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading a web page: %v\n", err)
		os.Exit(1)
	}
}
