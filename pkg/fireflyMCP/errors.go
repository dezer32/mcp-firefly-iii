package fireflyMCP

import "errors"

// Common errors for functional options
var (
	// Server option errors
	ErrNilHTTPClient    = errors.New("http client cannot be nil")
	ErrInvalidTimeout   = errors.New("timeout must be greater than 0")
	ErrEmptyAPIToken    = errors.New("API token cannot be empty")
	ErrEmptyBaseURL     = errors.New("base URL cannot be empty")
	ErrNilRequestEditor = errors.New("request editor cannot be nil")
	ErrNilMiddleware    = errors.New("middleware cannot be nil")
	ErrInvalidRateLimit = errors.New("rate limit and burst must be greater than 0")
	ErrNilConfig        = errors.New("config cannot be nil")
	
	// Client option errors
	ErrInvalidRetryCount = errors.New("retry count cannot be negative")
)