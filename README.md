# TicTacToe

This revised classic game of TicTacToe with Go/HTMX, is a real-time game where two players compete against each other over a network, and other observe the game in progress. Read more about implementation details here: link to article

![tic-tac-toe](./tictactoe/game.png)


NOTE: You can find the single-player [branch here](https://github.com/abayomipopoola/tictactoe/tree/single-player)

## Technology

The multiplayer function over a network uses the long-polling technique for communication between the client and the server. This decision was driven by two main factors: simplicity and the short-lived nature of the application. While WebSockets offer continuous real-time updates, setting up a WebSocket server and maintaining a persistent connection can be more complex and resource-intensive than a straightforward long-polling solution, particularly for applications that don't require constant real-time feedback.

However, one challenge associated with long-polling is the potential risk of missing game events during brief client disconnections. To address this, I integrated a simple message queue data-structure. This ensures that even if a client temporarily loses connection or misses a message, it can retrieve missed game events upon reconnection. Clients send the timestamp of their last received message to maintain synchronization, allowing the server to provide any updates that might have been missed during disconnections.

- Go - https://go.dev/
- HTMX - https://htmx.org/

### Server

The server is written in Go, not just for its simplicity and conciseness, but its builtin concurrency support.


### Client

HTMX///