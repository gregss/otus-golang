package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/config"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/logger"
	queue "github.com/gregss/otus/hw12_13_14_15_calendar/internal/queue/producer"
	sqlstorage "github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jasonlvhit/gocron"
)

func main() {
	time.Sleep(5 * time.Second)

	cnf := &config.Config{}
	config.LoadConfig(cnf)
	logg := logger.New(cnf.Logger.Level, cnf.Logger.File)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := sqlstorage.New(cnf.Storage.Dsn)

	go func() {
		<-ctx.Done()
		err := storage.Close(ctx)
		if err != nil {
			logg.Error(err.Error())
		}
	}()

	if err := storage.Connect(ctx); err != nil {
		logg.Error(err.Error())
	}

	if err := storage.Ping(); err != nil {
		logg.Error(err.Error())
	}

	calendar := app.New(logg, storage)

	err := gocron.Every(1).Minute().Do(func() {
		// кривой способ, по хорошему нужен еще метод в апп/сторадж получения событий по дате времени начала.
		for _, e := range calendar.DayEvents(ctx, time.Now()) {
			err := queue.Publish(
				cnf.Queue.URI,
				"calendar-exchange",
				"direct",
				"calendar",
				fmt.Sprintf("%s: %s - %s", e.Time, e.Title, e.Description),
			)
			if err != nil {
				logg.Error(err.Error())
			}
			log.Printf("published %dB OK", len("body"))
		}
	})
	if err != nil {
		logg.Error(err.Error())
	}

	err = gocron.Every(1).Day().Do(func() {
		err := calendar.DelPrevYearEvents(ctx)
		if err != nil {
			logg.Error(err.Error())
		}
	})

	if err != nil {
		logg.Error(err.Error())
	}

	select {}
}
