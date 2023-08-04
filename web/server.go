package web

import (
	"fmt"
	"net/http"
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
}

type Move struct {
	GamePlay
	CreatedAt int64
}

// Player Post request
type PlayRequest struct {
	Player Player `json:"player"`
	Row    int    `json:"row"`
	Column int    `json:"column"`
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

	// chi custom logger
	h.Use(middleware.Logger)
	h.Use(middleware.Timeout(30 * time.Second))

	// web handlers
	h.Get("/", h.Home())
	h.Route("/web", func(r chi.Router) {
		r.Get("/init", h.Init())
		r.Get("/updates", h.Updates())
		r.Get("/player/{playerID}/move", h.Move())
		r.Delete("/reset", h.Reset())
	})

	return h
}

func (h *Handler) Init() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if new player can join the game
		if ok, pID := h.pool.CanPlay(); ok {
			jsonString := fmt.Sprintf(`{"id": "%s", "ts": "%s"}`, pID, h.pool.play[pID])
			w.Write([]byte(jsonString))
			return
		}

		player := r.URL.Query().Get("player")
		ts := r.URL.Query().Get("ts")
		// check for active game
		if (player == "X" || player == "O") && h.pool.play[player] == ts {
			jsonString := fmt.Sprintf(`{"id": "%s", "ts": "%s"}`, player, ts)
			w.Write([]byte(jsonString))
			return
		}

		jsonString := fmt.Sprintf(`{"id": "%s", "ts": "%s"}`, "-1", "")
		w.Write([]byte(jsonString))
	}
}

func (h *Handler) Updates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lastUpdate := r.URL.Query().Get("lastUpdate")
		lastUpdateUnix, _ := strconv.ParseInt(lastUpdate, 10, 64)
		getMoves := func() []Move {
			moves := h.queue.Copy()
			filtered := []Move{}
			for _, move := range moves {
				if move.CreatedAt > lastUpdateUnix {
					filtered = append(filtered, move)
				}
			}
			return filtered
		}

		moves := getMoves()
		if len(moves) > 0 {
			// update the http client hypermedia rep of game state
			_ = Home(w, HomeParams{moves[0]})
			return
		}

		ch, close := h.pubsub.Subscribe()
		defer close()

		select {
		case <-ch:
			moves = getMoves()
			// update the http client hypermedia rep of game state
			_ = Home(w, HomeParams{moves[0]})
		case <-r.Context().Done():
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("request timed out"))
		}
	}
}
