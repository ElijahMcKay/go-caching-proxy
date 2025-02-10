package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ElijahMcKay/go-caching-proxy/proxy"
)

func ProxyHandler(proxyServer *proxy.ProxyServer) http.HandlerFunc {
	var cacheHeaderVal string
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		url, err := url.JoinPath(proxyServer.Origin, r.URL.Path)
		if err != nil {
			fmt.Println("Error creating request URL")
			http.Error(w, "Failed to create URL to origin server", http.StatusInternalServerError)
		}

		// Example cache key: "GET-http://dummyjson.com/products"
		cacheKey := strings.Join([]string{r.Method, url}, "-")

		// check if data exists and is not stale.  Continue on refreshing request if it's stale
		cachedValue, cacheExists, isStale := proxyServer.ReadCache(cacheKey)
		if cacheExists && !isStale {
			cacheHeaderVal = "HIT"
			proxyServer.WriteResponse(cacheHeaderVal, cachedValue, w, r)
			elapsed := time.Since(start)
			fmt.Printf("Cache %s latency: %v\n", cacheHeaderVal, elapsed)
			return
		}

		// making request to origin if data isn't found in cache
		response, err := http.Get(url)
		if err != nil {
			fmt.Print("Issue making request", err)
		}
		defer response.Body.Close()

		// writes to cache and returns the newly created entry
		cacheEntry, err := proxyServer.WriteCache(cacheKey, response)
		if err != nil {
			fmt.Printf("Error creating cache entry: %v", err)
		}

		cacheHeaderVal = "MISS"
		proxyServer.WriteResponse(cacheHeaderVal, cacheEntry, w, r)
		elapsed := time.Since(start)
		fmt.Printf("Cache %s latency: %v\n", cacheHeaderVal, elapsed)
	}
}
