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

func GenerateCommand(ctx context.Context, logger *log.Logger, cfg *GenerateConfig) error {
	if cfg.Depth <= 0 {
		return fmt.Errorf("depth cannot be null or negative")
	}

	logger.Printf("generating new map of size %d", cfg.Depth)

	var seed int64
	if cfg.Seed != "" {
		seed = int64(crc32.ChecksumIEEE([]byte(cfg.Seed)))
		logger.Printf("using seed %s", cfg.Seed)
	} else {
		seed = time.Now().UnixNano()
	}

	rand.Seed(int64(seed))

	cities := invader.NewCities()
	cities.GenerateRandomCity(cfg.Depth)
	cities.Print(os.Stdout)
	return nil
}

func generateCommand(ctx context.Context, logger *log.Logger, rcfg *RootConfig, args []string) *ffcli.Command {
	var cfg GenerateConfig
	cfg.RootConfig = rcfg

	flagSet := flag.NewFlagSet("generate", flag.ExitOnError)
	flagSet.StringVar(&cfg.Seed, "seed", "", "the seed used to generate the map, empty seed will be choose if empty")
	flagSet.IntVar(&cfg.Depth, "depth", 100, "the depth of the wanted map")

	return &ffcli.Command{
		Name:        "generate",
		ShortUsage:  "invader generate",
		ShortHelp:   "generate a new random cities with the given size",
		FlagSet:     flagSet,
		Subcommands: []*ffcli.Command{},
		Exec: func(ctx context.Context, args []string) error {
			return GenerateCommand(ctx, logger, &cfg)
		},
	}
}
