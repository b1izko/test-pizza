package main

import (
	"log"

	srv "github.com/b1izko/test-pizza/manager/internal"
	"github.com/b1izko/test-pizza/manager/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Port         = ":3500"
	DataBaseAddr = "mongodb://localhost:27017"
	DataBaseName = "pizza"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	svc := srv.New(&srv.Config{
		ListenAddress: Port,
		DBName:        DataBaseName,
		DBUrl:         DataBaseAddr,
	}, ch)

	if err := svc.Start(); err != nil {
		log.Fatalf("Cannot start service: %v", err)
	}

	defer svc.Close()
}
