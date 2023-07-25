package main

import (
	"fmt"
	"net/http"

	"github.com/abayomipopoola/game/tictactoe"
	"github.com/abayomipopoola/game/web"
)

func main() {
	// Create new tictactoe game
	t := tictactoe.NewTicTacToe()

	// Register http handlers
	h := web.NewHandler(t)

	fmt.Println("Game server started:::")
	http.ListenAndServe(":3000", h)
}
