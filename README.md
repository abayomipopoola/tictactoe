# TicTacToe

This revised classic game of TicTacToe with Go/HTMX, is a real-time game where two players compete against each other over a network, and other observe the game in progress. Read more about implementation details here: link to article

![tic-tac-toe](./game.png)


NOTE: You can find the single-player [branch here](https://github.com/abayomipopoola/tictactoe/tree/single-player)

## Technology
The game employs long-polling techniques for the communication bridge between the client and the server. When the server receives a message (moves), it broadcasts this to all active connections. A more robust solution might involve using WebSockets. An important consideration is the potential for events transpiring while clients are reconnecting; therefore, integrating a message queue becomes essential. Clients would send the timestamp of the last received message to maintain synchronization. This aspect hasn't been added to the current version.

- Go - https://go.dev/
- HTMX - https://htmx.org/

### Server

...


### Client

...