package tictactoe

import "fmt"

type Player string

const (
	X Player = "X"
	O Player = "O"

	BoardSize = 3
)

// Define Tic-Tac-Toe game interface
type Game interface {
	NewGame()
	Move(player Player, row, column int) error
	GetBoard() [BoardSize][BoardSize]*Player
	GetTurn() Player
	GetWinner() *Player
}

type TicTacToe struct {
	turn   Player
	winner *Player
	board  [BoardSize][BoardSize]*Player
}

// NewTicTacToe creates a new instance of a TicTacToe game
func NewTicTacToe() *TicTacToe {
	return &TicTacToe{
		turn:   X,
		winner: nil,
		board:  [BoardSize][BoardSize]*Player{},
	}
}

// NewGame reset its state to a fresh game
func (t *TicTacToe) NewGame() {
	t.turn = X
	t.winner = nil
	t.board = [BoardSize][BoardSize]*Player{}
}

// Move places the player's symbol in the given row and column and errors if move is invalid
func (t *TicTacToe) Move(player Player, row, column int) error {
	if t.winner != nil {
		return fmt.Errorf("tictactoe: game is already over")
	}
	if player != t.turn {
		return fmt.Errorf("tictactoe: not %s's turn", player)
	}
	if !isValidMove(t.board, row, column) {
		return fmt.Errorf("tictactoe: location %d,%d is not empty or out of bounds", row, column)
	}

	t.board[row][column] = &player

	t.winner = getWinner(t.board)
	if t.winner == nil {
		t.turn = switchPlayer(t.turn)
	}

	return nil
}

// GetTurn returns the turn of the TicTacToe game
func (t *TicTacToe) GetTurn() Player {
	return t.turn
}

// GetWinner returns the winner of the TicTacToe game
func (t *TicTacToe) GetWinner() *Player {
	return t.winner
}

// GetBoard returns the TicTacToe board
func (t *TicTacToe) GetBoard() [BoardSize][BoardSize]*Player {
	return t.board
}

func isValidMove(board [BoardSize][BoardSize]*Player, row, column int) bool {
	return row >= 0 && row < BoardSize && column >= 0 && column < BoardSize && board[row][column] == nil
}

func getWinner(board [BoardSize][BoardSize]*Player) *Player {
	winConditions := [][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {1, 2}},
		{{2, 0}, {2, 1}, {2, 2}},
		{{0, 0}, {1, 0}, {2, 0}},
		{{0, 1}, {1, 1}, {2, 1}},
		{{0, 2}, {1, 2}, {2, 2}},
		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
	}

	for _, line := range winConditions {
		a, b, c := line[0], line[1], line[2]
		if board[a[0]][a[1]] != nil && board[b[0]][b[1]] != nil && board[c[0]][c[1]] != nil {
			if *board[a[0]][a[1]] == *board[b[0]][b[1]] && *board[a[0]][a[1]] == *board[c[0]][c[1]] {
				return board[a[0]][a[1]]
			}
		}
	}

	return nil
}

func switchPlayer(currentPlayer Player) Player {
	if currentPlayer == X {
		return O
	}
	return X
}
