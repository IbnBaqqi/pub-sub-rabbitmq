package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ibnbaqqi/pub-sub-rabbitmq/internal/pubsub"
	"github.com/ibnbaqqi/pub-sub-rabbitmq/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"
	
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("unable to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	amqpCh, err := conn.Channel()
	if err != nil {
		fmt.Printf("unable to create channel: %v", err)
	}

	err = pubsub.PublishJSON(
		amqpCh,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: true},
	)
	if err != nil {
		log.Printf("could not publish time: %v", err)
	}
	fmt.Println("Pause message sent!")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("RabbitMQ connection closed.")
}
