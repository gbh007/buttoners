package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
	generator "github.com/gbh007/buttoners/tools/gg/internal"
)

type Config struct {
	Boards []struct {
		UID           string `toml:"uid"`
		Name          string `toml:"name"`
		OutputFile    string `toml:"output_file"`
		ServiceFilter string `toml:"service_filter"`
	} `toml:"boards"`
}

func main() {
	configPath := flag.String("config", "gg-config.toml", "path to config")

	flag.Parse()

	var cfg Config
	logger := slog.Default()

	f, err := os.Open(*configPath)
	if err != nil {
		logger.Error("open config", "error", err)
		os.Exit(1)
	}

	_, err = toml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		logger.Error("decode config", "error", err)
		os.Exit(1)
	}

	for i, board := range cfg.Boards {

		if board.UID == "" {
			logger.Error("empty board uid", "number", i)
		}

		if board.OutputFile == "" {
			logger.Error("empty board output file", "number", i)
		}

		g := generator.New(board.UID, board.Name, board.ServiceFilter)

		dashboardModel, err := g.Build()
		if err != nil {
			panic(err)
		}

		out, err := os.Create(board.OutputFile)
		if err != nil {
			panic(err)
		}

		enc := json.NewEncoder(out)
		enc.SetIndent("", "   ")

		err = enc.Encode(dashboardModel)
		if err != nil {
			panic(err)
		}

		err = out.Close()
		if err != nil {
			panic(err)
		}
	}
}
