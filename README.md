# TicTacToe

This revised classic game of TicTacToe with Go/HTMX, is a real-time game where two players compete against each other over a network, and other observe the game in progress. Read more about implementation details here: link to article

![tic-tac-toe](./game.png)

NOTE: You can find the single-player [branch here](https://github.com/abayomipopoola/tictactoe/tree/single-player)

## Technology

The multiplayer function over a network uses the long-polling technique for communication between the client and the server. This decision was driven by two main factors: simplicity and the short-lived nature of the application. While WebSockets offer continuous real-time updates, setting up a WebSocket server and maintaining a persistent connection can be more complex and resource-intensive than a straightforward long-polling solution, particularly for applications that don't require constant real-time feedback.

However, one challenge associated with long-polling is the potential risk of missing game events during brief client disconnections. To address this, I integrated a simple message queue data-structure. This ensures that even if a client temporarily loses connection or misses a message, it can retrieve missed game events upon reconnection. Clients send the timestamp of their last received message to maintain synchronization, allowing the server to provide any updates that might have been missed during disconnections.

### Server

The server is written in [Go](https://go.dev), not just for its simplicity and conciseness, but also its built-in concurrency and a robust standard library. Making it easier to implement simple long-polling server with a message queue logic to handle missing game events in case of connection drop -- without using third-party linraries.

The server is written in Go, not only for its simplicity and conciseness but also because of its native concurrency support and a comprehensive standard library. These features make it straightforward to implement a long-polling server with a message queue mechanism. This ensures uninterrupted gameplay, even if a connection drops—all without relying on third-party libraries.

For a detailed code walkthrough, refer to this [tutorial](https://github.com/abayomipopoola/tictactoe/tree/single-player)

### Client

[HTMX](https://htmx.org) is the main reason for this simple game. HTMX is a fascinating project because it allows developers to access modern web features directly from HTML, sidestepping the complexity of traditional JavaScript frameworks. This makes web interactivity and AJAX integrations simpler, lighter, and more maintainable. 

Using HTMX can result in faster development, improved performance, and a more direct connection between your server-side code and the front-end.

For a detailed code walkthrough, refer to this [tutorial](https://github.com/abayomipopoola/tictactoe/tree/single-player)