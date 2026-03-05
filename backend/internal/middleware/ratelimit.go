package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"pineapple-music/internal/util"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*ipLimiter
	rps      rate.Limit
	burst    int
}

func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*ipLimiter),
		rps:      rate.Limit(rps),
		burst:    burst,
	}

	// Evict stale entries every 3 minutes
	go rl.cleanupLoop()

	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rps, rl.burst)
		rl.visitors[ip] = &ipLimiter{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting for static resources (covers, streams)
		path := c.Request.URL.Path
		if strings.HasSuffix(path, "/cover") || strings.HasSuffix(path, "/stream") {
			c.Next()
			return
		}

		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			util.ErrorResponse(c, http.StatusTooManyRequests, "rate_limited", "Too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
