package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	. "github.com/abayomipopoola/game/tictactoe"
	"github.com/stretchr/testify/assert"
)

// TestMove tests the Move handler function
func TestMove(t *testing.T) {
	game := NewTicTacToe()
	handler := NewHandler(game)

	playRequest := PlayRequest{Player: X, Row: 0, Column: 0}
	reqBody, _ := json.Marshal(playRequest)

	req, err := http.NewRequest("POST", "/game/move", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGame(t *testing.T) {
	// Create a new game and handler
	game := NewTicTacToe()
	handler := NewHandler(game)

	// Make a few moves
	game.Move(X, 0, 0)
	game.Move(O, 1, 1)
	game.Move(X, 2, 2)

	// Make a GET request to the game endpoint
	req := httptest.NewRequest("GET", "/game", nil)
	rec := httptest.NewRecorder()
	handler.Game().ServeHTTP(rec, req)

	// Check that the response status is OK
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the response
	var gp GamePlay
	if err := json.Unmarshal(rec.Body.Bytes(), &gp); err != nil {
		t.Errorf("unable to parse response body: %v", err)
	}

	// Validate that the response matches the game state
	if !reflect.DeepEqual(gp.Board, game.GetBoard()) {
		t.Errorf("Expected board to be %+v, but got %+v", game.GetBoard(), gp.Board)
	}

	if gp.PlayerTurn != game.GetTurn() {
		t.Errorf("Expected turn to be %v, but got %v", game.GetTurn(), gp.PlayerTurn)
	}

	winner := game.GetWinner()
	if winner == nil && gp.Winner != nil || winner != nil && gp.Winner == nil {
		t.Errorf("Expected winner to be %v, but got %v", winner, gp.Winner)
	}
}
