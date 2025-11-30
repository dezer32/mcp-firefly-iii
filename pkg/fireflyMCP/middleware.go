package fireflyMCP

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// FireflyTokenKey is the context key for storing the Firefly III API token
const FireflyTokenKey contextKey = "firefly_token"

// GetTokenFromContext extracts the Firefly III API token from the request context
func GetTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(FireflyTokenKey).(string); ok {
		return token
	}
	return ""
}

// TokenExtractionMiddleware creates middleware that extracts the Firefly III API token
// from the Authorization header and stores it in the request context.
// The token is then used by the MCP server to authenticate with Firefly III API.
// Health check endpoints (/health, /ready) are excluded from token requirement.
func TokenExtractionMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip token extraction for health check endpoints
			if r.URL.Path == "/health" || r.URL.Path == "/ready" {
				next.ServeHTTP(w, r)
				return
			}

			auth := r.Header.Get("Authorization")
			if auth == "" {
				logger.Warn("missing authorization header",
					"remote_addr", r.RemoteAddr,
					"path", r.URL.Path)
				http.Error(w, "Authorization: Bearer <firefly-token> required", http.StatusUnauthorized)
				return
			}

			// Check Bearer prefix (case-insensitive)
			if len(auth) < 7 || !strings.EqualFold(auth[:7], "bearer ") {
				logger.Warn("invalid authorization format",
					"remote_addr", r.RemoteAddr,
					"path", r.URL.Path)
				http.Error(w, "Authorization: Bearer <firefly-token> required", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(auth[7:])
			if token == "" {
				logger.Warn("empty token",
					"remote_addr", r.RemoteAddr,
					"path", r.URL.Path)
				http.Error(w, "Empty token", http.StatusUnauthorized)
				return
			}

			// Store token in context for use by MCP handlers
			ctx := context.WithValue(r.Context(), FireflyTokenKey, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CORSMiddleware creates middleware for Cross-Origin Resource Sharing.
func CORSMiddleware(allowedOrigins []string, logger *slog.Logger) func(http.Handler) http.Handler {
	allowAll := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"

	originSet := make(map[string]bool)
	for _, origin := range allowedOrigins {
		originSet[origin] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			if allowAll {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if origin != "" && originSet[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			// Set other CORS headers
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Mcp-Session-Id, Mcp-Protocol-Version")
			w.Header().Set("Access-Control-Expose-Headers", "Mcp-Session-Id, Mcp-Protocol-Version")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware creates per-IP rate limiting middleware.
type rateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	b        int
}

func newRateLimiter(r float64, b int) *rateLimiter {
	return &rateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        rate.Limit(r),
		b:        b,
	}
}

func (rl *rateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.limiters[ip] = limiter
	}

	return limiter
}

func RateLimitMiddleware(requestsPerSecond float64, burst int, logger *slog.Logger) func(http.Handler) http.Handler {
	limiter := newRateLimiter(requestsPerSecond, burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip rate limiting for health checks
			if r.URL.Path == "/health" || r.URL.Path == "/ready" {
				next.ServeHTTP(w, r)
				return
			}

			ip := getClientIP(r)
			if !limiter.getLimiter(ip).Allow() {
				logger.Warn("rate limit exceeded",
					"ip", ip,
					"path", r.URL.Path)
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP from the request, handling proxies.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestLoggingMiddleware creates middleware for structured request logging.
func RequestLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status
			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			// Log request details
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"remote_addr", getClientIP(r),
				"user_agent", r.UserAgent())
		})
	}
}
