package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	ff "github.com/peterbourgon/ff/v3"
	ffcli "github.com/peterbourgon/ff/v3/ffcli"
)

var (
	debugMode bool
	logger    *log.Logger
)

type RootConfig struct {
	Debug bool
}

var rootFlagSet = flag.NewFlagSet("invader", flag.ExitOnError)

func parseRootConfig(args []string) (*RootConfig, error) {
	var cfg RootConfig

	rootFlagSet.BoolVar(&cfg.Debug, "debug", false, "enable debug mode")

	if err := ff.Parse(rootFlagSet, args); err != nil {
		return nil, fmt.Errorf("unable to parse flags: %w", err)
	}

	return &cfg, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logWriter := io.Discard

	args := os.Args[1:]
	rcfg, err := parseRootConfig(args)
	if err != nil {
		fmt.Println("error:", err.Error())
		os.Exit(1)
	}

	if rcfg.Debug {
		logWriter = os.Stderr
	}

	logger := log.New(logWriter, "debug: ", log.LstdFlags)

	root := &ffcli.Command{
		Name:    "invader [flags] <subcommand>",
		FlagSet: rootFlagSet,
		Exec:    func(_ context.Context, _ []string) error { return flag.ErrHelp },
		Subcommands: []*ffcli.Command{
			startCommand(ctx, logger, rcfg, args),
			generateCommand(ctx, logger, rcfg, args),
		},
	}

	if err := root.ParseAndRun(ctx, args); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		}
		os.Exit(1)
	}
}
