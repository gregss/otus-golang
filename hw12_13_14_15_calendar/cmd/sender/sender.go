package main

import (
	"log"
	"time"

	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/config"
	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/logger"
	queue "github.com/gregss/otus/hw12_13_14_15_calendar/internal/queue/consumer"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/streadway/amqp"
)

func main() {
	time.Sleep(5 * time.Second)

	cnf := &config.Config{}
	config.LoadConfig(cnf)

	logg := logger.New(cnf.Logger.Level, cnf.Logger.File)

	c, err := queue.NewConsumer(
		cnf.Queue.URI,
		"calendar-exchange",
		"direct",
		"calendar-queue",
		"calendar",
		"calendar-consumer",
		handle,
	)
	if err != nil {
		logg.Error(err.Error())
	}

	var forever chan struct{}
	<-forever

	log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		logg.Error(err.Error())
	}
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
