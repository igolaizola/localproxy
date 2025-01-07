package localproxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/elazarl/goproxy"
)

// Run starts a local proxy server that forwards traffic to an upstream server.
// It returns the URL of the local proxy server.
//   - The ctx is used to detect when the process should be shut down.
//   - The debug parameter enables verbose logging.
//   - The upstream parameter is the URL of the upstream server.
//   - The addr parameter is the listen address of the local proxy server, leave
//     it empty to listen on all interfaces and use a random port.
func Run(ctx context.Context, debug bool, upstream, addr string) (string, error) {
	if upstream == "" {
		return "", fmt.Errorf("localproxy: upstream URL is required")
	}
	if addr == "" {
		addr = "0.0.0.0:0"
	}

	// Parse the upstream URL
	parsed, err := url.Parse(upstream)
	if err != nil {
		return "", fmt.Errorf("localproxy: couldn't parse upstream URL: %w", err)
	}

	// Create a goproxy instance
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = debug

	// Configure the proxy to forward traffic to the upstream
	proxy.Tr = &http.Transport{
		Proxy: http.ProxyURL(parsed),
	}

	// Listen on the provided listen addr (e.g., "127.0.0.1:0" or "0.0.0.0:0")
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("localproxy: couldn't listen on %s: %w", addr, err)
	}

	// Extract the host and port
	host, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		return "", fmt.Errorf("localproxy: couldn't split host and port: %w", err)
	}
	// If the host is 0.0.0.0 or ::, replace it with 127.0.0.1
	if host == "0.0.0.0" || host == "[::]" || host == "::" {
		host = "127.0.0.1"
	}

	// Construct the final local proxy URL
	localProxyURL := fmt.Sprintf("http://%s:%s", host, port)

	// Create an HTTP server using the goproxy as the handler
	server := &http.Server{
		Handler: proxy,
	}

	// Run the server in a background goroutine
	go func() {
		// If Serve() exits with a non-ErrServerClosed error, log it
		if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("localproxy: server error: %v\n", err)
		}
	}()

	// Watch for ctx cancellation to gracefully shut down the server
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	return localProxyURL, nil
}
