package middlewares

import (
	"net"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

type LimiterMiddleware struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
}

func NewLimiterMiddleware() *LimiterMiddleware {
	return &LimiterMiddleware{
		limiters: make(map[string]*rate.Limiter),
	}
}

func (l *LimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := l.getIP(r)
		if !l.isValidIP(ip) {
			http.Error(w, "Can't get client IP", http.StatusInternalServerError)
			return
		}

		limiter := l.getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests for IP "+ip, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (l *LimiterMiddleware) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	if limiter, ok := l.limiters[ip]; ok {
		return limiter
	}

	l.limiters[ip] = rate.NewLimiter(rate.Limit(30), 1)

	return l.limiters[ip]
}

func (l *LimiterMiddleware) getIP(r *http.Request) string {
	for _, header := range []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"CF-Connecting-IP", // Cloudflare
		"True-Client-IP",   // Akamai
	} {
		if value := r.Header.Get(header); value != "" {
			// X-Forwarded-For may contaings few IP. Take first
			if header == "X-Forwarded-For" {
				ips := strings.Split(value, ",")
				if len(ips) > 0 {
					return strings.TrimSpace(ips[0])
				}
			}
			return strings.TrimSpace(value)
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (l *LimiterMiddleware) isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}
