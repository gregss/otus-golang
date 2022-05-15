package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/config"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/gregss/otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/gregss/otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	time.Sleep(5 * time.Second)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	cnf := &config.Config{}
	config.LoadConfig(cnf)
	logg := logger.New(cnf.Logger.Level, cnf.Logger.File)

	var storage app.Storage
	if cnf.Storage.Type == "memory" {
		storage = memorystorage.New()
	} else {
		s := sqlstorage.New(cnf.Storage.Dsn)
		fmt.Printf("connect to (%v)", cnf.Storage.Dsn)

		go func() {
			<-ctx.Done()
			_ = s.Close(ctx)
		}()

		if err := s.Connect(ctx); err != nil {
			fmt.Printf("error connect storage (%v)", err)
			return
		}

		if err := s.Ping(); err != nil {
			fmt.Printf("error ping storage (%v)", err)
			return
		}

		storage = s
	}
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(*logg, *calendar, cnf.Server.Hport)
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)

			logg.Info("http server running")
		}
	}()

	internalgrpc.StartServer(*logg, *calendar, cnf.Server.Gport)
}
