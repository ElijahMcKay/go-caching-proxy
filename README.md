# go-caching-proxy
A simple locally hosted proxy server written in Go that caches HTTP requests in memory. Written as a submission for https://roadmap.sh/projects/caching-server

## Usage
Build the binary:
```go build -o capr main.go```
Note: Call the binary whatever you want.  `capr` seemed like a good abbreviate of "caching proxy"

View CLI flags:
```./capr --help```

Output:
```
Usage of ./capr:
  -origin string
        The remote server the proxy server will forward requests to (default "https://dummyjson.com")
  -port int
        The port to run the proxy server on (default 3000)
  -ttl int
        The duration in minutes before cache entries will be marked as stale. The next time a stale cache entry is read, it will be discarded and refreshed (default 10)
```

Start the server:
```./capr --port 3005 --ttl 5 --origin https://dummyjson.com```

Output:
```
Starting caching server at http://localhost:3000
Proxying requests to https://dummyjson.com
Press c + Enter to clear cache
```

Note: Press c + Enter in the terminal the proxy server is running in to clear all cache

