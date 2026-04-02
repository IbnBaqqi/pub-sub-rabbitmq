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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}
	queueName := routing.PauseKey + "." + username

	gameState := gamelogic.NewGameState(username)
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.TransientQueue,
		handlerPause(gameState),
	)
	if err != nil {
		log.Fatalf("Client unable to subscribe: %v", err)
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "move":
			armyMove, moveErr := gameState.CommandMove(words)
			if moveErr != nil {
				log.Printf("error occured: %v", moveErr)
				continue
			}
			// TODO publish the move
			log.Printf("Unit moved to %s", armyMove.ToLocation)
		case "spawn":
			spawnErr := gameState.CommandSpawn(words)
			if spawnErr != nil {
				log.Printf("error occured: %v", spawnErr)
				continue
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			// TODO: publish n malicious logs
			log.Printf("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			log.Printf("unknown command")
		}
	}
}

// handlerPause returns a handler that updates the game state when a pause/resume
// message is received
func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {

	return func(ps routing.PlayingState) {
        defer fmt.Print("> ")
        gs.HandlePause(ps)
    }
}