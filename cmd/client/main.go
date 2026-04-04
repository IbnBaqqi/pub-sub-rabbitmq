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

	fmt.Println("Starting Peril client...")

	// client connection
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("unable to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}
	gs := gamelogic.NewGameState(username)

	if err = pubsub.SubscribeJSON( // Subscribe to handle gamestate pause
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username, // queue name: "pause.bob"
		routing.PauseKey,              // routing key: matches "pause"
		pubsub.TransientQueue,
		handlerPause(gs),
	); err != nil {
		log.Fatalf("Client unable to subscribe: %v", err)
	}

	if err = pubsub.SubscribeJSON( // Subscribe to player army move
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+username, // queue name: "army_moves.bob"
		routing.ArmyMovesPrefix+".*",         // binding key: matches any "army_moves.<word>"
		pubsub.TransientQueue,
		handlerMove(gs),
	); err != nil {
		log.Fatalf("Client unable to subscribe: %v", err)
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "move":
			armyMove, err := gs.CommandMove(words)
			if err != nil {
				log.Printf("error: %v\n", err)
				continue
			}
			err = pubsub.PublishJSON( // publish the move
				publishCh,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+".*",
				armyMove,
			)
			if err != nil {
				log.Printf("error: %v\n", err)
				continue
			}
			fmt.Printf("Moved %v units to %s", len(armyMove.Units), armyMove.ToLocation)
		case "spawn":
			err := gs.CommandSpawn(words)
			if err != nil {
				log.Printf("error occured: %v", err)
				continue
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			log.Printf("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			log.Printf("unknown command")
		}
	}
}
