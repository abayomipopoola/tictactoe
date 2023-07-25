# TicTacToe

Classic game of TicTacToe with Go/HTMX. 

![tic-tac-toe](./game.png)

- Go - https://go.dev/
- HTMX - https://htmx.org/

## Endpoints

There are three main endpoints:

1. `GET /game`: This endpoint returns the current game state including the current player's turn, the game board state, and the winner if there is one.
    
2. `POST /game/move`: This endpoint is used to make a move. The player, row, and column should be provided in the request body in JSON format. For example:
    
```json
{ 
	"player": "X",   
	"row": 1,   
	"column": 2 
}
```

 On a successful move, this endpoint returns the updated game state in JSON. If the move is not successful, it returns a descriptive error message in JSON.
    
3. `DELETE /game`: This endpoint resets the game. The game board is cleared and the next GET request to `/game` will return a new game with an empty board.
    
## Running the Server

> Make sure you have Go installed (https://go.dev/doc/install).

Run the command below from the home directory:

`$ go run main.go`

This will start the game server on port 3000.

## Play game on your web browser

1. Access the home page at `http://localhost:3000/`.
2. To make a move, click on any unoccupied square within the displayed game board.
3. Click the "Reset" button located below the game board to restart the game.


## Play via REST 

You can play the game using a REST client like Postman or `cURL`. 

1. To get the current game state:

```bash
curl http://localhost:3000/game
```

2. To make a move (for example, player X moves to row 0 and column 0):

```bash
curl -X POST -H "Content-Type: application/json" -d '{"player":"X","row":0,"column":0}' http://localhost:3000/game/move
```
    
3. To reset the game:

```bash
curl -X DELETE http://localhost:3000/game
```