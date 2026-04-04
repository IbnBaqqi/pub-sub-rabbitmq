package main

import (
	"fmt"

	"github.com/ibnbaqqi/pub-sub-rabbitmq/internal/gamelogic"
	"github.com/ibnbaqqi/pub-sub-rabbitmq/internal/routing"
)

// handlerPause returns a handler that updates the game state when a pause/resume
// message is received
func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {

	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}

// handlerMove returns a handler that processes incoming army move messages
// from other players and updates the local game state accordingly.
func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {

	return func(armyMove gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(armyMove)
	}
}