package web

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/abayomipopoola/game/polling"
	. "github.com/abayomipopoola/game/tictactoe"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Active game play representation
type GamePlay struct {
	Board      [BoardSize][BoardSize]*Player
	PlayerTurn Player
	Winner     *Player
	Counter    int
}

type Move struct {
	GamePlay
	CreatedAt int64
}

type Handler struct {
	*chi.Mux
	game   Game
	pool   *PlayerPool
	queue  *polling.Queue[Move]
	pubsub *polling.PubSub
}

func NewHandler(game Game) *Handler {
	h := &Handler{
		Mux:    chi.NewMux(),
		game:   game,
		pool:   NewPlayerPool(2),
		queue:  polling.NewQueue[Move](9),
		pubsub: polling.NewPubSub(),
	}

	// chi custom middlewares
	h.Use(middleware.Logger)
	h.Use(middleware.Timeout(30 * time.Second))

	// web handlers
	h.Get("/", h.Home())
	h.Route("/web", func(r chi.Router) {
		r.Get("/updates", h.Updates())
		r.Get("/player/{playerID}/move", h.Move())
		r.Delete("/reset", h.Reset())
	})

	return h
}

func (h *Handler) Updates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getMovesAfter := func(timestamp int64) []Move {
			moves := h.queue.Copy()
			index := sort.Search(len(moves), func(i int) bool {
				return moves[i].CreatedAt > timestamp
			})
			return moves[index:]
		}

		lastUpdate, _ := strconv.ParseInt(r.URL.Query().Get("lastUpdate"), 10, 64)
		moves := getMovesAfter(lastUpdate)

		if len(moves) == 0 {
			ch, close := h.pubsub.Subscribe()
			defer close()

			select {
			case <-ch:
				moves = getMovesAfter(lastUpdate)
			case <-r.Context().Done():
				w.WriteHeader(http.StatusRequestTimeout)
				w.Write([]byte("request timed out"))
				return
			}
		}

		// If there are any moves, update the client.
		_ = renderHomePage(w, HomeParams{moves[0]})
	}
}
