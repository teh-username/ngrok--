## ngrok--

An attempt on an inferior implementation of [ngrok](https://github.com/inconshreveable/ngrok).

## Scope and Limitations

We'll only be considering a singular tunnel of HTTP traffic.

## How ngrok-- works

[ngrok](https://github.com/inconshreveable/ngrok/blob/master/docs/DEVELOPMENT.md) provides a succint explanation of how it works on high level which we'll be using as our reference.

For ngrok-- though, we'll be cutting corners a bit but the core technique of bypassing Firewall/NAT will be implemented, albeit differently. It's the thought that counts right?

Traffic flow should resemble the "diagram" below:

```
                              +--------------------------------+
+----+        +------+        | +------+        +------------+ |
|user|<------>|server|<------>| |client|<------>|local server| |
+----+        +------+        | +------+        +------------+ |
                              +--------------------------------+
                               Firewall/NAT
```

## Resources

* https://github.com/inconshreveable/ngrok/blob/master/docs/DEVELOPMENT.md
* https://www.stavros.io/posts/proxying-two-connections-go/
* https://github.com/nu7hatch/areyoufuckingcoding.me/blob/master/content/2012/08/03/tricks-tips-adapt-blocking-io-to-channel.md
* https://blog.cloudflare.com/this-is-strictly-a-violation-of-the-tcp-specification/
