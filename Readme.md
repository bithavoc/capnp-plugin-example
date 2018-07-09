# capnproto external plugin example

An example of how [go-capnproto2](https://github.com/capnproto/go-capnproto2) can be used to launch an external process and communicate via Stdin and Stdout(plugin-like).

## Compiling

```
(cd server; go build)
```

```
(cd client; go build)
```

## Run Server

```
(cd server; ./server)
```

Open server/debug.log, you should see an output similar to this:

```
debug2018/07/08 19:29:30 Debug started
debug2018/07/08 19:29:30 Creating connection
debug2018/07/08 19:29:30 connection open
debug2018/07/08 19:29:30 hash client open
debug2018/07/08 19:29:30 hello written
debug2018/07/08 19:29:30 world written
debug2018/07/08 19:29:30 will now call Sum
debug2018/07/08 19:29:30 sha1: 0a0a9f2a6772942557ab5355d76af442f8f65e01
```

## Re-generate cap and proto Go code:
```
capnp compile -I~/go/src/zombiezen.com/go/capnproto2/std -ogo hashes/hashes.capnp
```

**Note:**
The `three-way-plugin` branch contains an unsuccessful attempt to create a three way plugin, meaning a plugin that consumes RPC objects from another plugin, but since go-capnproto2 doesn't support Level 3, it does not work.