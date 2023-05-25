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

// StartCommand begins the simulation of the alien invasion.
func StartCommand(ctx context.Context, logger *log.Logger, cfg *StartConfig) error {
	var err error

	reader := os.Stdin
	if cfg.File != "" {
		if reader, err = os.Open(cfg.File); err != nil {
			return fmt.Errorf("unable to open file `%s`: %w", cfg.File, err)
		}

		logger.Printf("Reading `%s` file map", cfg.File)
	}

	ai := invader.NewAlienInvaders(logger, os.Stdout)

	if err = ai.ParseMap(reader); err != nil {
		return fmt.Errorf("unable parse the given map: %w", err)
	}

	if err := ai.GenerateAliens(cfg.NAlien); err != nil {
		return fmt.Errorf("unable generate `%d` alien: %w", cfg.NAlien, err)
	}

	// Run simulation
	fmt.Printf("* Starting the simulation with %d aliens\n", cfg.NAlien)
	err = ai.Run(ctx, cfg.StepLimit)
	switch err {
	case nil: // Reached steps limit
		fmt.Printf("All aliens are exhausted after performing more than %d steps!\n", cfg.StepLimit)
	case invader.ErrAllAliensAreKO:
		fmt.Printf("All aliens have been killed/trapped!\n")
	default:
		return err
	}

	logger.Print("Simulation completed!")

	// Print the final state of the map.
	fmt.Printf("* final map:\n")
	ai.PrintMap()

	return nil
}

func startCommand(ctx context.Context, logger *log.Logger, rcfg *RootConfig, args []string) *ffcli.Command {
	var cfg StartConfig
	cfg.RootConfig = rcfg

	flagSet := flag.NewFlagSet("start", flag.ExitOnError)
	flagSet.IntVar(&cfg.NAlien, "aliens", 4, "The number of aliens that will be generated on the map")
	flagSet.IntVar(&cfg.StepLimit, "max_steps", 10000, "The maximum number of steps an alien can perform before becoming exhausted.")
	flagSet.StringVar(&cfg.File, "file", "", "Read from a specified file instead of the standard input.")

	return &ffcli.Command{
		Name:       "start",
		ShortUsage: "invader start -alien [value] -file [path] -max_steps [value]",
		ShortHelp:  "Start the invader simulation by reading from the standard input.",
		LongHelp: `This subcommand initiates the Alien Invaders simulation. The
program reads from standard input by default, but you can
specify a file insteaad.`,
		FlagSet:     flagSet,
		Subcommands: []*ffcli.Command{},
		Exec: func(ctx context.Context, args []string) error {
			return StartCommand(ctx, logger, &cfg)
		},
	}
}
