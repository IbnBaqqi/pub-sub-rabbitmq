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

	// publish server channel
	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("unable to create channel: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,    // exchange
		routing.GameLogSlug,           // queue name
		routing.GameLogSlug + ".*",    // key
		pubsub.DurableQueue,
	)
	if err != nil {
		log.Fatalf("unable to declare queue: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	// print server commands
	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			publishToExchange(publishCh, true)
		case "resume":
			publishToExchange(publishCh, false)
		case "quit":
			fmt.Println("RabbitMQ connection closed.")
			return
		default:
			log.Println("unknown command")
		}
	}
}

// publishToExchange publishes a pause or resume state to the peril direct exchange.
// Set isPause to true to pause the game, false to resume.
func publishToExchange(publishCh *amqp.Channel, isPause bool) {

	err := pubsub.PublishJSON(
		publishCh,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: isPause},
	)
	if err != nil {
		log.Printf("could not publish time: %v", err)
	}
	fmt.Println("Pause message published!")
}
