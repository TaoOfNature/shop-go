package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/TaoOfNature/shop-go/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type visitor struct {
	tokens     int
	lastRefill time.Time
}

func RateLimit(rps int, burst int) gin.HandlerFunc {
	var (
		mu       sync.Mutex
		visitors = make(map[string]*visitor)
	)

	refillEvery := time.Second
	if rps <= 0 {
		rps = 10
	}
	if burst <= 0 {
		burst = 20
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		v, ok := visitors[ip]
		if !ok {
			v = &visitor{tokens: burst, lastRefill: now}
			visitors[ip] = v
		}

		elapsed := now.Sub(v.lastRefill)
		if elapsed >= refillEvery {
			refill := int(elapsed/refillEvery) * rps
			v.tokens += refill
			if v.tokens > burst {
				v.tokens = burst
			}
			v.lastRefill = now
		}

		if v.tokens <= 0 {
			mu.Unlock()
			response.Error(c, http.StatusTooManyRequests, 1004, "too many requests")
			c.Abort()
			return
		}
		v.tokens--
		mu.Unlock()

		c.Next()
	}
}
