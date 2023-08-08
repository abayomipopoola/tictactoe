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

func Home(w io.Writer, p HomeParams) error {
	home := template.Must(template.ParseFS(content, "home.html"))
	return home.Execute(w, p)
}

func (h *Handler) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if new player can join the game
		if ok, c := h.pool.CanPlay(); ok {
			gamePlay := GamePlay{
				h.game.GetBoard(), h.game.GetTurn(), h.game.GetWinner(), c,
			}
			move := Move{gamePlay, time.Now().Unix()}
			h.queue.Enqueue(move)
			h.pubsub.Publish()
			_ = Home(w, HomeParams{Move{gamePlay, time.Now().Unix()}})
			return
		}
		gamePlay := GamePlay{
			h.game.GetBoard(), h.game.GetTurn(), h.game.GetWinner(), -1,
		}
		_ = Home(w, HomeParams{Move{gamePlay, time.Now().Unix()}})
	}
}

func (h *Handler) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.game.NewGame()
		h.pool = NewPlayerPool(2)
		count = 0
		gamePlay := GamePlay{
			h.game.GetBoard(), h.game.GetTurn(), h.game.GetWinner(), count,
		}

		move := Move{gamePlay, time.Now().Unix()}
		h.queue.Enqueue(move)
		h.pubsub.Publish()
		_ = Home(w, HomeParams{move})
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

		gamePlay := GamePlay{
			h.game.GetBoard(), h.game.GetTurn(), h.game.GetWinner(), count,
		}

		move := Move{gamePlay, time.Now().Unix()}
		h.queue.Enqueue(move)
		h.pubsub.Publish()
		_ = Home(w, HomeParams{move})
	}
}
