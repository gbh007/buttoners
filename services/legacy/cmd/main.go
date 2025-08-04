package main

import (
	"github.com/gbh007/buttoners/services/legacy/internal/controller"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("addr", ":8080", "web server address")
	debug := flag.Bool("d", false, "debug mode")
	dbType := flag.String("db", "sqlite", "db type sqlite, postgres, mysql")
	conn := flag.String("conn", "test.db", "db connection string")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	ll := slog.LevelInfo
	if *debug {
		ll = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: ll}))

	c, err := controller.New(
		logger,
		*addr,
		*debug,
		*dbType,
		*conn,
	)
	if err != nil {
		logger.Error("create controller", "error", err)
		os.Exit(1)
	}

	logger.Info("start server")

	err = c.Serve(ctx)
	if err != nil {
		logger.Error("serve http", "error", err)
		os.Exit(1)
	}

	logger.Info("have a nice day")
}
