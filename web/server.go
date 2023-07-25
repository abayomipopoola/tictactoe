package web

import (
	"encoding/json"
	"net/http"

	. "github.com/abayomipopoola/game/tictactoe"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Define Tic-Tac-Toe game interface
type Game interface {
	NewGame()
	Move(player Player, row, column int) error
	GetBoard() [BoardSize][BoardSize]*Player
	GetTurn() Player
	GetWinner() *Player
}

// Active game play representation
type GamePlay struct {
	Board      [BoardSize][BoardSize]*Player
	PlayerTurn Player
	Winner     *Player
}

// Player Post request
type PlayRequest struct {
	Player Player `json:"player"`
	Row    int    `json:"row"`
	Column int    `json:"column"`
}

type Handler struct {
	*chi.Mux
	game Game
}

func NewHandler(game Game) *Handler {
	h := &Handler{
		Mux:  chi.NewMux(),
		game: game,
	}

	// chi custom logger
	h.Use(middleware.Logger)

	// REST endpoints for game logic
	h.Route("/game", func(r chi.Router) {
		r.Get("/", h.Game())      // GET /game
		r.Post("/move", h.Move()) // POST /game/move
		r.Delete("/", h.Reset())  // DELETE /game
	})

	// Extras: HTML web handlers
	h.Get("/", h.Home())
	h.Route("/web", func(r chi.Router) {
		r.Get("/player/{playerID}/move", h.MoveWeb())
		r.Delete("/reset", h.ResetWeb())
	})

	return h
}

func (h *Handler) Game() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSONResponse(w, h.game)
	}
}

func (h *Handler) Move() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var play PlayRequest

		// decode request body into person struct
		err := json.NewDecoder(r.Body).Decode(&play)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.game.Move(play.Player, play.Row, play.Column)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJSONResponse(w, h.game)
	}
}

func (h *Handler) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.game.NewGame()

		// return a 204 No Content status
		w.WriteHeader(http.StatusNoContent)
	}
}

// utility function to write JSON response
func writeJSONResponse(w http.ResponseWriter, g Game) {
	gamePlay := GamePlay{
		g.GetBoard(),
		g.GetTurn(),
		g.GetWinner(),
	}
	// set content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// use json package to marshal GamePlay
	jsonData, err := json.Marshal(gamePlay)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write JSON response
	w.Write(jsonData)
}
