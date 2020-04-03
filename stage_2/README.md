## Stage 2

We'll be shifting gears on this stage by implementing the following flow:
_Note: We're assuming that the control connection of the previous stage has already been established_

1) [Server] Open a public and proxy port
2) [Server] On new public connection, send command to create proxy to client over the control connection
3) [Client] Dial a tcp connection to proxy port of server
4) [Server] On new proxy connection, pipe traffic between public and proxy connection

or visually:

```
+----+         +------+         +------+
|user|-------->|server|         |client|
+----+         +------+         +------+
                  ^                 ^
                  |                 |
                  |                 |
                  +-----------------+
                      create_proxy
                      command

------------------------------------------

                       proxy
                       connection
                  +-----------------+
                  |                 |
                  |                 |
                  v                 v
+----+         +------+         +------+
|user|-------->|server|         |client|
+----+         +------+         +------+
                  ^                 ^
                  |                 |
                  |                 |
                  +-----------------+
```

Now we see the importance of the control connection. It basically gives the server the capability of on-demand proxy tunnel creation whenever a public user attempts to access the "localhost" content.

The client can then just pipe the traffic between the proxy connection and the localhost server, which we'll tackle on the next stage.

### Verification

To verify that everything is working, run `netcat localhost 60624`. Anything you write should be echoed back to you. Type "stop" to finish the session.

You can also run `watch -n 1 ss -nt dst 127.0.0.1` on a separate terminal to see the tcp sockets being utilized.

[Previous](../stage_1/README.md)
