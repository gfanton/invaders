package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gfanton/invader"
	ffcli "github.com/peterbourgon/ff/v3/ffcli"
)

type StartConfig struct {
	*RootConfig

	StepLimit int
	NAlien    int
	File      string
}

func StartCommand(ctx context.Context, logger *log.Logger, cfg *StartConfig) error {
	var err error
	logger.Printf("start with %d aliens", cfg.NAlien)

	reader := os.Stdin
	if cfg.File != "" {
		if reader, err = os.Open(cfg.File); err != nil {
			return fmt.Errorf("unable to open file `%s`: %w", cfg.File, err)
		}

		logger.Printf("reading `%s` file map", cfg.File)
	}

	ai := invader.NewAlienInvaders(logger)

	if err = ai.ParseMap(reader); err != nil {
		return fmt.Errorf("unable parse map: %w", err)
	}

	if err := ai.GenerateAliens(cfg.NAlien); err != nil {
		return fmt.Errorf("unable generate `%d` alien: %w", cfg.NAlien, err)
	}

	// run simulation
	err = ai.Run(ctx, cfg.StepLimit)
	switch err {
	case nil: // reached steps limit
		fmt.Printf("all aliens are exhausted by doing more than %d steps!\n", cfg.StepLimit)
	case invader.ErrAllAliensAreDead:
		fmt.Printf("all aliens have been killed!\n")
	default:
		return err
	}

	logger.Print("simulation is done!")

	ai.PrintMap()

	return nil
}

func startCommand(ctx context.Context, logger *log.Logger, rcfg *RootConfig, args []string) *ffcli.Command {
	var cfg StartConfig
	cfg.RootConfig = rcfg

	flagSet := flag.NewFlagSet("start", flag.ExitOnError)
	flagSet.IntVar(&cfg.NAlien, "alien", 10, "the seed used to start the map, empty seed will be choose if empty")
	flagSet.IntVar(&cfg.StepLimit, "max_steps", 10000, "the maximum step aliens can do before beeing exhausted")
	flagSet.StringVar(&cfg.File, "file", "", "read the target file instead of stdin")

	return &ffcli.Command{
		Name:        "start",
		ShortUsage:  "invader start",
		ShortHelp:   "start invader simulation reading stdin",
		FlagSet:     flagSet,
		Subcommands: []*ffcli.Command{},
		Exec: func(ctx context.Context, args []string) error {
			return StartCommand(ctx, logger, &cfg)
		},
	}
}
