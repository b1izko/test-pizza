package main

import (
	"github.com/b1izko/test-pizza/internal/logger"
	"github.com/b1izko/test-pizza/manager/internal/config"
	srv "github.com/b1izko/test-pizza/manager/internal/server"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conf = config.New().Load()
)

func main() {

	conn, err := amqp.Dial(conf.RabbitMQ.URL)
	if logger.IsError(err, "Failed to connect to RabbitMQ") {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if logger.IsError(err, "Failed to open a channel") {
		panic(err)
	}
	defer ch.Close()

	conf = config.New().Load()
	svc := srv.New(&srv.Config{
		ListenAddress: conf.Port,
		DBName:        conf.MongoDB.Name,
		DBUrl:         conf.MongoDB.URL,
	}, ch)

	err = svc.Start()
	if logger.IsError(err, "Cannot start service") {
		panic(err)
	}

	defer svc.Close()
}
