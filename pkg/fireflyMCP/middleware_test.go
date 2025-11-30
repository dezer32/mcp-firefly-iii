package fireflyMCP

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestBearerAuthMiddleware_ValidToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := BearerAuthMiddleware("secret-token", testLogger())
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestBearerAuthMiddleware_InvalidToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := BearerAuthMiddleware("secret-token", testLogger())
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestBearerAuthMiddleware_MissingHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := BearerAuthMiddleware("secret-token", testLogger())
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestBearerAuthMiddleware_NoAuthRequired(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Empty token means no auth required
	middleware := BearerAuthMiddleware("", testLogger())
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestBearerAuthMiddleware_HealthEndpointBypass(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := BearerAuthMiddleware("secret-token", testLogger())
	wrapped := middleware(handler)

	// Health endpoint should bypass auth
	for _, path := range []string{"/health", "/ready"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()

		wrapped.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "path: %s", path)
	}
}

func TestBearerAuthMiddleware_CaseInsensitive(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := BearerAuthMiddleware("secret-token", testLogger())
	wrapped := middleware(handler)

	// Test case-insensitive "Bearer" prefix
	for _, authHeader := range []string{"Bearer secret-token", "bearer secret-token", "BEARER secret-token"} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", authHeader)
		rr := httptest.NewRecorder()

		wrapped.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "auth: %s", authHeader)
	}
}

func TestCORSMiddleware_AllowAll(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CORSMiddleware([]string{"*"}, testLogger())
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CORSMiddleware([]string{"http://allowed.com", "http://also-allowed.com"}, testLogger())
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://allowed.com")
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, "http://allowed.com", rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_DisallowedOrigin(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CORSMiddleware([]string{"http://allowed.com"}, testLogger())
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://not-allowed.com")
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_Preflight(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	middleware := CORSMiddleware([]string{"*"}, testLogger())
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Contains(t, rr.Header().Get("Access-Control-Allow-Methods"), "POST")
}

func TestRateLimitMiddleware_WithinLimit(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Allow 10 requests per second with burst of 5
	middleware := RateLimitMiddleware(10, 5, testLogger())
	wrapped := middleware(handler)

	// First 5 requests should succeed (burst)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rr := httptest.NewRecorder()

		wrapped.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "request %d", i)
	}
}

func TestRateLimitMiddleware_ExceedsLimit(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Very low limit: 0.1 requests per second with burst of 1
	middleware := RateLimitMiddleware(0.1, 1, testLogger())
	wrapped := middleware(handler)

	// First request succeeds
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.RemoteAddr = "192.168.1.2:12345"
	rr1 := httptest.NewRecorder()
	wrapped.ServeHTTP(rr1, req1)
	assert.Equal(t, http.StatusOK, rr1.Code)

	// Second request should be rate limited
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.RemoteAddr = "192.168.1.2:12345"
	rr2 := httptest.NewRecorder()
	wrapped.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusTooManyRequests, rr2.Code)
}

func TestRateLimitMiddleware_DifferentIPs(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Very low limit: 0.1 requests per second with burst of 1
	middleware := RateLimitMiddleware(0.1, 1, testLogger())
	wrapped := middleware(handler)

	// Request from IP 1
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.RemoteAddr = "192.168.1.100:12345"
	rr1 := httptest.NewRecorder()
	wrapped.ServeHTTP(rr1, req1)
	assert.Equal(t, http.StatusOK, rr1.Code)

	// Request from IP 2 should also succeed (different rate limit bucket)
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.RemoteAddr = "192.168.1.200:12345"
	rr2 := httptest.NewRecorder()
	wrapped.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusOK, rr2.Code)
}

func TestRateLimitMiddleware_HealthBypass(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Very low limit
	middleware := RateLimitMiddleware(0.001, 1, testLogger())
	wrapped := middleware(handler)

	// Health endpoints should bypass rate limiting
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.RemoteAddr = "192.168.1.3:12345"
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	}
}

func TestGetClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1, 192.168.1.1")
	req.RemoteAddr = "127.0.0.1:12345"

	ip := getClientIP(req)

	assert.Equal(t, "10.0.0.1", ip)
}

func TestGetClientIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "10.0.0.2")
	req.RemoteAddr = "127.0.0.1:12345"

	ip := getClientIP(req)

	assert.Equal(t, "10.0.0.2", ip)
}

func TestGetClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.100:54321"

	ip := getClientIP(req)

	assert.Equal(t, "192.168.1.100", ip)
}

func TestRequestLoggingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestLoggingMiddleware(testLogger())
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
