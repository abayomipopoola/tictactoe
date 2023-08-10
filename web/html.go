package web

import (
	"embed"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	. "github.com/abayomipopoola/game/tictactoe"
	"github.com/go-chi/chi/v5"
)

//go:embed home.html
var content embed.FS

var count int

type HomeParams struct {
	Move
}

// Handlers:

func (h *Handler) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var move Move
		if ok, c := h.pool.CanPlay(); ok {
			gamePlay := h.createGamePlay(c)
			move = h.enqueueAndPublish(gamePlay)
		} else {
			gamePlay := h.createGamePlay(-1)
			move = Move{gamePlay, time.Now().Unix()}
		}
		_ = renderHomePage(w, HomeParams{move})
	}
}

func (h *Handler) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.game.NewGame()
		h.pool = NewPlayerPool(2)
		count = 0
		gamePlay := h.createGamePlay(count)
		move := h.enqueueAndPublish(gamePlay)
		_ = renderHomePage(w, HomeParams{move})
	}
}

func (h *Handler) Move() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		player := chi.URLParam(r, "playerID")
		row, _ := strconv.Atoi(r.URL.Query().Get("row"))
		col, _ := strconv.Atoi(r.URL.Query().Get("col"))
		count++

		err := h.game.Move(Player(player), row, col)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		gamePlay := h.createGamePlay(count)
		move := h.enqueueAndPublish(gamePlay)
		_ = renderHomePage(w, HomeParams{move})
	}
}

// Helpers:

func renderHomePage(w io.Writer, p HomeParams) error {
	home := template.Must(template.ParseFS(content, "home.html"))
	return home.Execute(w, p)
}

func (h *Handler) createGamePlay(c int) GamePlay {
	return GamePlay{
		Board:      h.game.GetBoard(),
		PlayerTurn: h.game.GetTurn(),
		Winner:     h.game.GetWinner(),
		Counter:    c,
	}
}

func (h *Handler) enqueueAndPublish(gamePlay GamePlay) Move {
	move := Move{gamePlay, time.Now().Unix()}
	h.queue.Enqueue(move)
	h.pubsub.Publish()
	return move
}
