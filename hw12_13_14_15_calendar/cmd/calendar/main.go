package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/gregss/otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/gregss/otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/jackc/pgx/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level, config.Logger.File)

	var storage app.Storage
	if config.Storage.Type == "memory" {
		storage = memorystorage.New()
	} else {
		s := sqlstorage.New(config.Storage.Dsn)

		go func() {
			<-ctx.Done()
			_ = s.Close(ctx)
		}()

		if err := s.Connect(ctx); err != nil {
			fmt.Printf("%v", err)
		}

		if err := s.Ping(); err != nil {
			fmt.Printf("%v", err)
		}

		storage = s
	}
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(*logg, *calendar, config.Server.Hport)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	internalgrpc.StartServer(*logg, *calendar, config.Server.Gport)

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) // nolint:gocritic
	}
}
