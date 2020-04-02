## Stage 1

In this stage, we'll implement a simple client-side proxy that can tunnel tcp packets to and from a target local server. We'll be iterating over this code in later stages.

It should look like:

```
 +----+
 |curl|------+
 +----+      v
          +------+       +-------------------+
          |client|------>|local/remote server|
          +------+       +-------------------+
             ^
+-------+    |
|browser|----+
+-------+
```

### Verification

To verify, run the client code `go run main.go` then execute `curl localhost:54286` on a separate console. You can also try to open `http://localhost:54286` on your browser.

[Home](../README.md)
