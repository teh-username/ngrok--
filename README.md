## ngrok--

An iterative attempt on an inferior implementation of [ngrok](https://github.com/inconshreveable/ngrok).

### Scope and Limitations

We'll only be considering a singular tunnel of HTTP traffic.

### Stages

* [Stage 0](stage_0/README.md) - Proof of concept implementation of a simple client-side proxy
* [Stage 1](stage_1/README.md) - Base server and client implementation plus intro to the control connection

### How ngrok-- works

[ngrok](https://github.com/inconshreveable/ngrok/blob/master/docs/DEVELOPMENT.md) provides a succint explanation of how it works on high level which we'll be using as our reference.

For ngrok-- though, we'll be cutting corners a bit but the core technique of bypassing Firewall/NAT will be implemented, albeit differently. It's the thought that counts right?

Traffic flow should resemble the "diagram" below (WIP):

```
                              +--------------------------------+
+----+        +------+        | +------+        +------------+ |
|user|<------>|server|<------>| |client|<------>|local server| |
+----+        +------+        | +------+        +------------+ |
                              +--------------------------------+
                               Firewall/NAT
```
