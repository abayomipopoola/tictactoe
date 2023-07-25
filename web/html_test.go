package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/abayomipopoola/game/tictactoe"
)

func TestMoveWeb(t *testing.T) {
	// Create a new game and handler
	game := NewTicTacToe()
	handler := NewHandler(game)

	// The player, row, and col to test
	playerID := "X"
	row := 0
	col := 0

	req := httptest.NewRequest("GET", fmt.Sprintf("/web/player/%s/move?row=%d&col=%d", playerID, row, col), nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Check that the response status is OK
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that the move was made in the game
	if game.GetBoard()[row][col] == nil || *game.GetBoard()[row][col] != Player(playerID) {
		t.Errorf("expected move to be made at row %d col %d", row, col)
	}

	// Check that the turn has changed
	if game.GetTurn() == Player(playerID) {
		t.Errorf("expected turn to change but it's still %s", playerID)
	}
}
