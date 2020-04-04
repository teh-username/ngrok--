## Stage 1

In this stage, we'll be creating our base server and client that does the following:

1. Server comes online and opens a port for our control connection.
2. Client comes online and establishes the control connection to the server.
3. Server accepts the connection initiated by Client and waits for traffic.
4. Client sends the initial test traffic "ping". Server then replies with "pong".
5. Repeat (4) 10 times before the client terminates the connection.

### [Aside] Control Connection
The client-initiated control connection can be considered to be the foundation of how ngrok works. Using this long-lived connection, the server can now freely orchestrate the client to setup the seemingly impossible tunnel to localhost. We'll tackle more about this "orchestration" in later stages.

In [ngrok's case](https://github.com/inconshreveable/ngrok/blob/master/docs/DEVELOPMENT.md#wire-format), it uses netstrings as the format for the commands send down the control connection. For ngrok--, we'll just go with newline ('\n') delimited commands.


[Previous](../stage_0/README.md) <<>> [Next](../stage_2/README.md)
