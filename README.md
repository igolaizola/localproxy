# localproxy

A lightweight HTTP proxy that forwards traffic to an upstream server, built on top of [goproxy](https://github.com/elazarl/goproxy).  
Use it as a command-line tool or embed it in your Go application to run a local proxy server that directs all requests to your chosen upstream.

> 📢 Connect with us! Join our Telegram group for support and collaboration: [t.me/igohub](https://t.me/igohub)

## 📦 Installation

### Command-line tool

You can install the **localproxy** CLI via Go:

```bash
go install github.com/igolaizola/localproxy/cmd/localproxy@latest
```

Or you can download the binary from the [releases](https://github.com/igolaizola/localproxy/releases)

## 🛠️ Usage

### As a command-line tool

Run `localproxy --help` to view available flags and subcommands:

```bash
localproxy --help
```

Basic usage:

```bash
localproxy --upstream https://example.com --addr 0.0.0.0:8080
```

**Flags:**

- `--upstream` (required): The URL of the upstream server to forward traffic to.
- `--addr`: The local address and port to listen on (defaults to `0.0.0.0:0`).
- `--debug`: Enable verbose logging for debugging purposes.
- `--config`: Path to a config file (YAML by default).
- `--help`: Show usage and exit.

> By default, if you leave `--addr` empty (e.g. `0.0.0.0:0`), **localproxy** picks a random available port.

**Version command:**

```bash
localproxy version
```

This displays version information, such as the semantic version, commit hash, and build date (if available).

### Configuration file

You can store all the flags in a YAML file instead of passing them in the CLI:

```bash
localproxy --config localproxy.yaml
```

An example `localproxy.yaml`:

```yaml
upstream: https://example.com
addr: 0.0.0.0:8080
debug: true
```

**Environment variables**  
You can also override flags using environment variables by prefixing them with `LOCALPROXY_`. For example:

```bash
export LOCALPROXY_UPSTREAM="https://example.com"
export LOCALPROXY_DEBUG="true"
localproxy
```

## 📚 Using localproxy as a Go library

You can integrate **localproxy** into your own Go projects. Just import and call `Run`:

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/igolaizola/localproxy"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    localProxyURL, err := localproxy.Run(ctx, true, "https://example.com", "127.0.0.1:8080")
    if err != nil {
        log.Fatalf("Failed to start localproxy: %v", err)
    }

    log.Printf("Proxy running at %s", localProxyURL)

    // Wait for the context to expire or manual cancellation
    <-ctx.Done()
    log.Println("Shutting down localproxy")
}
```

- **Parameters:**
  - `ctx`: a `context.Context` used to shut down the proxy gracefully.
  - `debug`: a boolean to enable or disable verbose logs.
  - `upstream`: the URL of the upstream server.
  - `addr`: the listen address for your proxy server.

The function returns a URL that you can use to configure your clients (e.g., set `HTTP_PROXY` or `HTTPS_PROXY` environment variables to point to the localproxy).

## 💖 Support

If you find **localproxy** helpful, please give the repository a star ⭐!

You can also support the project financially:

Invite me for a coffee at ko-fi (0% fees):

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/igolaizola)

Or at buymeacoffee:

[![buymeacoffee](https://user-images.githubusercontent.com/11333576/223217083-123c2c53-6ab8-4ea8-a2c8-c6cb5d08e8d2.png)](https://buymeacoffee.com/igolaizola)

Donate to my PayPal:

[paypal.me/igolaizola](https://www.paypal.me/igolaizola)

Sponsor me on GitHub:

[github.com/sponsors/igolaizola](https://github.com/sponsors/igolaizola)

Thank you for your support!

## 🔗 Resources

- [elazarl/goproxy](https://github.com/elazarl/goproxy) — The HTTP proxy library used internally.
