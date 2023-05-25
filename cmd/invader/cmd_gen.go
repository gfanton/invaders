package main

import (
	"context"
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"math/rand"
	"os"
	"time"

	invader "github.com/gfanton/invader"
	ffcli "github.com/peterbourgon/ff/v3/ffcli"
)

type GenerateConfig struct {
	*RootConfig

	Depth int
	Seed  string
}

// GenerateCommand generates a new random city map and prints it to stdout.
func GenerateCommand(ctx context.Context, logger *log.Logger, cfg *GenerateConfig) error {
	if cfg.Depth <= 0 {
		return fmt.Errorf("depth cannot be null or negative")
	}

	logger.Printf("generating new map of size %d", cfg.Depth)

	if cfg.Seed == "" {
		cfg.Seed = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	seed := int64(crc32.ChecksumIEEE([]byte(cfg.Seed)))
	logger.Printf("using seed %s", cfg.Seed)

	// set seed
	rand.Seed(seed)

	cities := invader.NewCities()
	cities.GenerateRandomCity(cfg.Depth)
	cities.Print(os.Stdout)
	return nil
}

func generateCommand(ctx context.Context, logger *log.Logger, rcfg *RootConfig, args []string) *ffcli.Command {
	var cfg GenerateConfig
	cfg.RootConfig = rcfg

	flagSet := flag.NewFlagSet("generate", flag.ExitOnError)
	flagSet.StringVar(&cfg.Seed, "seed", "", "the seed used to generate the map; a random seed will be chosen if left empty")
	flagSet.IntVar(&cfg.Depth, "depth", 5, "the depth of the desired map")

	return &ffcli.Command{
		Name:        "generate",
		ShortUsage:  "invader generate -depth [value] -seed [string]",
		ShortHelp:   "generate a new random cities with the given depth",
		LongHelp:    "This subcommand is used to generate a new random city map of a given depth.",
		FlagSet:     flagSet,
		Subcommands: []*ffcli.Command{},
		Exec: func(ctx context.Context, args []string) error {
			return GenerateCommand(ctx, logger, &cfg)
		},
	}
}
