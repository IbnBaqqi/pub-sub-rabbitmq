package main

import (
	"fmt"
	"log"

	"github.com/ibnbaqqi/pub-sub-rabbitmq/internal/gamelogic"
	"github.com/ibnbaqqi/pub-sub-rabbitmq/internal/pubsub"
	"github.com/ibnbaqqi/pub-sub-rabbitmq/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	fmt.Println("Starting Peril server...")

	// sever connection
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("unable to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	// server connection
	amqpCh, err := conn.Channel()
	if err != nil {
		fmt.Printf("unable to create channel: %v", err)
	}

	// print server commands
	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) < 1 {
			continue
		}
		if words[0] == "pause" {
			publishToExchange(amqpCh, true)
		}
		if words[0] == "resume" {
			publishToExchange(amqpCh, false)
		}
		if words[0] == "quit" {
			fmt.Println("RabbitMQ connection closed.")
			break
		}
	}
}

func publishToExchange(amqpCh *amqp.Channel, isPause bool) {

	err := pubsub.PublishJSON(
		amqpCh,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: isPause},
	)
	if err != nil {
		log.Printf("could not publish time: %v", err)
	}
	fmt.Println("Pause message sent!")
}
