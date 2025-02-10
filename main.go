package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ElijahMcKay/go-caching-proxy/handlers"
	"github.com/ElijahMcKay/go-caching-proxy/proxy"
)

func main() {
	// defining flags for the CLI
	port := flag.Int("port", 3000, "The port to run the proxy server on")
	origin := flag.String("origin", "https://dummyjson.com", "The remote server the proxy server will forward requests to")
	cacheEntryTTL := flag.Int("ttl", 10, "The duration in minutes before cache entries will be marked as stale. The next time a stale cache entry is read, it will be discarded and refreshed")
	flag.Parse()

	// if flag.NFlag() == 0 {
	// 	// Print help if no arguments are passed
	// 	flag.PrintDefaults()
	// 	return
	// }
	// create Proxy server
	proxyServer := proxy.NewProxy(*origin, *cacheEntryTTL)

	go startWebServer(proxyServer, *port)
	go captureTerminalInput(proxyServer)

	select {}

}

func startWebServer(proxyServer *proxy.ProxyServer, port int) {

	// setup handler for all requests
	http.Handle("/", handlers.ProxyHandler(proxyServer))

	// startup server
	fmt.Printf("Caching server listening at http://localhost:%v\nProxying requests to %s\n\n", port, proxyServer.Origin)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func captureTerminalInput(proxyServer *proxy.ProxyServer) {
	time.Sleep(1 * time.Second)
	// create scanner to read input from terminal
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Press c + Enter to clear cache\n")
		if scanner.Scan() {
			input := scanner.Text()
			if strings.TrimSpace(input) == "c" {
				proxyServer.ClearCache()
			} else {
				fmt.Println("Invalid input:", scanner.Err())
				break
			}
		}
	}
}
