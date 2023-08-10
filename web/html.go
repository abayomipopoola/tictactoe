package web

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"

	. "github.com/abayomipopoola/game/tictactoe"
	"github.com/go-chi/chi/v5"
)

//go:embed home.html
var content embed.FS

type HomeParams struct {
	Game GamePlay
}

// Handlers:

func (h *Handler) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gamePlay := h.createGamePlay()
		_ = renderHomePage(w, HomeParams{gamePlay})
	}
}

func (h *Handler) ResetWeb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.game.NewGame()
		gamePlay := h.createGamePlay()
		_ = renderHomePage(w, HomeParams{gamePlay})
	}
}

func (h *Handler) MoveWeb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		player := chi.URLParam(r, "playerID")
		row, _ := strconv.Atoi(r.URL.Query().Get("row"))
		col, _ := strconv.Atoi(r.URL.Query().Get("col"))

		err := h.game.Move(Player(player), row, col)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		gamePlay := h.createGamePlay()
		_ = renderHomePage(w, HomeParams{gamePlay})
	}
}

// Helpers:

func (h *Handler) createGamePlay() GamePlay {
	return GamePlay{
		Board:      h.game.GetBoard(),
		PlayerTurn: h.game.GetTurn(),
		Winner:     h.game.GetWinner(),
	}
}

func renderHomePage(w io.Writer, p HomeParams) error {
	home := template.Must(template.ParseFS(content, "home.html"))
	return home.Execute(w, p)
}

func writeJSONResponse(w http.ResponseWriter, g Game) {
	gamePlay := GamePlay{
		g.GetBoard(),
		g.GetTurn(),
		g.GetWinner(),
	}
	// set content type to application/json
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(gamePlay)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
