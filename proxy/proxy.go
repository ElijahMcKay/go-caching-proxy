package proxy

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ElijahMcKay/go-caching-proxy/cache"
)

type ProxyServer struct {
	Origin        string
	Cache         map[string]*cache.CacheObject
	CacheEntryTTL int
	mu            sync.RWMutex
}

func (p *ProxyServer) WriteCache(cacheKey string, r *http.Response) (*cache.CacheObject, error) {
	// http.Response.Body is a stream that is consumed, so we must store it for resuse
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body")
	}

	p.mu.Lock()
	p.Cache[cacheKey] = &cache.CacheObject{
		Response:     r,
		ResponseBody: body,
		CreatedAt:    time.Now(),
		TTL:          time.Duration(p.CacheEntryTTL) * time.Minute,
	}
	p.mu.Unlock()

	return p.Cache[cacheKey], nil
}

func (p *ProxyServer) ReadCache(cacheKey string) (*cache.CacheObject, bool, bool) {
	p.mu.RLock()
	cachedValue, cacheEntryExists := p.Cache[cacheKey]
	p.mu.RUnlock()

	isStale := false
	if cacheEntryExists {
		isStale = time.Since(cachedValue.CreatedAt) > cachedValue.TTL
	}

	return cachedValue, cacheEntryExists, isStale
}

func (p *ProxyServer) WriteResponse(cacheHeaderVal string, cacheEntry *cache.CacheObject, w http.ResponseWriter, r *http.Request) {

	for key, value := range cacheEntry.Response.Header {
		w.Header()[key] = value
	}
	w.Header().Add("X-Cache", cacheHeaderVal)
	w.Write(cacheEntry.ResponseBody)

}

func (p *ProxyServer) ClearCache() {
	p.Cache = make(map[string]*cache.CacheObject)
	fmt.Println("All cache cleared")
}

func NewProxy(origin string, cacheEntryTTL int) *ProxyServer {
	return &ProxyServer{
		Origin:        origin,
		Cache:         make(map[string]*cache.CacheObject),
		CacheEntryTTL: cacheEntryTTL,
		mu:            sync.RWMutex{},
	}
}
