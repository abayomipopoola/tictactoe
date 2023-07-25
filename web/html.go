package web

import (
	"embed"
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

func Home(w io.Writer, p HomeParams) error {
	home := template.Must(template.ParseFS(content, "home.html"))
	return home.Execute(w, p)
}

/**
*	The handler functions below are extras for playing the game on the web.
 */

func (h *Handler) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gamePlay := GamePlay{
			h.game.GetBoard(),
			h.game.GetTurn(),
			h.game.GetWinner(),
		}
		_ = Home(w, HomeParams{gamePlay})
	}
}

func (h *Handler) ResetWeb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.game.NewGame()
		gamePlay := GamePlay{
			h.game.GetBoard(),
			h.game.GetTurn(),
			h.game.GetWinner(),
		}
		_ = Home(w, HomeParams{gamePlay})
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
		gamePlay := GamePlay{
			h.game.GetBoard(),
			h.game.GetTurn(),
			h.game.GetWinner(),
		}
		_ = Home(w, HomeParams{gamePlay})
	}
}
