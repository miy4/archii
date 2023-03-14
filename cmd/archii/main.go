package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/miy4/archii"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type Setting struct {
	DocumentDir string `toml:"doc_dir"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}

	setting, err := readSettings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading a configuration file: %v\n", err)
		os.Exit(1)
	}

	fs := flag.NewFlagSet("archii", flag.ExitOnError)
	var dir string
	fs.StringVar(&dir, "d", ".", "output directory")
	fs.VisitAll(func(f *flag.Flag) {
		if f.Name == "d" && setting.DocumentDir != "" {
			f.Value.Set(setting.DocumentDir)
		}
	})

	root := &ffcli.Command{
		ShortUsage: "archii [-d dir] <URL>",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires an URL to read")
			}

			return archii.RunApp(args[0], dir)
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading a web page: %v\n", err)
		os.Exit(1)
	}
}

func readSettings() (Setting, error) {
	paths := []string{
		os.ExpandEnv("${HOME}/.config/archii/archii.toml"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			continue
		}

		var setting Setting
		_, err := toml.DecodeFile(path, &setting)
		if err != nil {
			return Setting{}, err
		}
		return setting, nil
	}

	return Setting{}, nil
}
