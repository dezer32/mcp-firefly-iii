# Stage 1: Builder
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install ca-certificates for HTTPS requests during build
RUN apk add --no-cache ca-certificates

# Download dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mcp-server ./cmd/mcp-server

# Stage 2: Runtime
FROM alpine:3.21

# Install ca-certificates for HTTPS requests to Firefly III
# Create non-root user for security
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 mcp

USER mcp

COPY --from=builder /build/mcp-server /usr/local/bin/

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/usr/local/bin/mcp-server", "--transport", "http"]
