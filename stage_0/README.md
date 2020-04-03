## Stage 0

In this stage, we'll implement a simple client-side proxy that can tunnel tcp packets to and from a target local server. We'll be using this code in later stages.

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

### References

* [Developer's guide to ngrok](https://github.com/inconshreveable/ngrok/blob/master/docs/DEVELOPMENT.md)
* [Proxying two connections in Go](https://www.stavros.io/posts/proxying-two-connections-go/)
* [GoTip #1 - Adapt blocking IO to a channel](https://github.com/nu7hatch/areyoufuckingcoding.me/blob/master/content/2012/08/03/tricks-tips-adapt-blocking-io-to-channel.md)
* [This is strictly a violation of the TCP specification](https://blog.cloudflare.com/this-is-strictly-a-violation-of-the-tcp-specification/)

[Home](../README.md) || [Next](../stage_1/README.md)
