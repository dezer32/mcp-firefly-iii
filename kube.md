# Kubernetes Deployment Guide

This guide covers deploying the Firefly III MCP Server with HTTP transport in Kubernetes.

## Prerequisites

- Kubernetes cluster (1.19+)
- kubectl configured
- Docker image built and pushed to a registry

## Building and Pushing the Image

```bash
# Build the image
docker build -t your-registry/firefly-mcp:latest .

# Push to registry
docker push your-registry/firefly-mcp:latest
```

## Configuration Options

All configuration is done via environment variables with the `FIREFLY_MCP_` prefix:

| Variable | Description | Default |
|----------|-------------|---------|
| `FIREFLY_MCP_SERVER_URL` | Firefly III API base URL | **Required** |
| `FIREFLY_MCP_LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `FIREFLY_MCP_CLIENT_TIMEOUT` | HTTP client timeout for Firefly III API (seconds) | `60` |
| `FIREFLY_MCP_HTTP_ENABLED` | Enable HTTP transport | `false` |
| `FIREFLY_MCP_HTTP_PORT` | HTTP server port | `8080` |
| `FIREFLY_MCP_HTTP_HOST` | HTTP server bind address | `0.0.0.0` |
| `FIREFLY_MCP_HTTP_READ_TIMEOUT` | HTTP read timeout (seconds) | `30` |
| `FIREFLY_MCP_HTTP_WRITE_TIMEOUT` | HTTP write timeout (seconds) | `30` |
| `FIREFLY_MCP_HTTP_IDLE_TIMEOUT` | HTTP idle timeout (seconds) | `120` |
| `FIREFLY_MCP_HTTP_SESSION_TIMEOUT` | MCP session timeout (seconds) | `300` |
| `FIREFLY_MCP_HTTP_RATE_LIMIT` | Requests per second per IP | `10.0` |
| `FIREFLY_MCP_HTTP_RATE_BURST` | Rate limit burst capacity | `20` |
| `FIREFLY_MCP_HTTP_ALLOWED_ORIGINS` | CORS allowed origins (comma-separated or `*`) | `*` |

## CLI Flags

The server supports command-line flags that override config file and environment variables:

```bash
# Run with HTTP transport
./mcp-server --transport=http --port=8080

# Available flags:
#   --transport   Transport mode: "stdio" (default) or "http"
#   --port        Override HTTP port
#   --config      Path to config file (default: config.yaml)
#   --log-level   Log level: debug, info, warn, error
```

## Kubernetes Manifests

### 1. Namespace (Optional)

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: firefly-mcp
```

### 2. Authentication

In HTTP mode, the MCP server does **not** store any Firefly III credentials. Instead, each client passes their own Firefly III Personal Access Token via the `Authorization` header:

```
Authorization: Bearer <firefly-iii-personal-access-token>
```

This enables multi-tenant usage where different clients can use their own Firefly III accounts.

**No secrets needed for Firefly III API tokens in Kubernetes!**

### 3. ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: firefly-mcp-config
  namespace: firefly-mcp
data:
  # Required: Firefly III API URL
  server-url: "https://firefly.example.com/api"
  # Log level: debug, info, warn, error
  log-level: "info"
  # HTTP client timeout for Firefly III API calls (seconds)
  client-timeout: "30"
  # HTTP Server Configuration
  http-read-timeout: "30"
  http-write-timeout: "30"
  http-idle-timeout: "120"
  http-session-timeout: "300"
  # Rate Limiting (per IP)
  http-rate-limit: "10.0"
  http-rate-burst: "20"
  # CORS (use "*" for all origins, or comma-separated list)
  http-allowed-origins: "*"
```

### 4. Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: firefly-mcp
  namespace: firefly-mcp
  labels:
    app: firefly-mcp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: firefly-mcp
  template:
    metadata:
      labels:
        app: firefly-mcp
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - name: firefly-mcp
        image: your-registry/firefly-mcp:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
        env:
        # Required: Firefly III server URL
        - name: FIREFLY_MCP_SERVER_URL
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: server-url
        # Log level
        - name: FIREFLY_MCP_LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: log-level
              optional: true
        # HTTP client timeout
        - name: FIREFLY_MCP_CLIENT_TIMEOUT
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: client-timeout
              optional: true
        # HTTP Server - Enable and bind
        - name: FIREFLY_MCP_HTTP_ENABLED
          value: "true"
        - name: FIREFLY_MCP_HTTP_PORT
          value: "8080"
        - name: FIREFLY_MCP_HTTP_HOST
          value: "0.0.0.0"
        # HTTP Server - Timeouts
        - name: FIREFLY_MCP_HTTP_READ_TIMEOUT
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: http-read-timeout
              optional: true
        - name: FIREFLY_MCP_HTTP_WRITE_TIMEOUT
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: http-write-timeout
              optional: true
        - name: FIREFLY_MCP_HTTP_IDLE_TIMEOUT
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: http-idle-timeout
              optional: true
        - name: FIREFLY_MCP_HTTP_SESSION_TIMEOUT
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: http-session-timeout
              optional: true
        # Rate Limiting
        - name: FIREFLY_MCP_HTTP_RATE_LIMIT
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: http-rate-limit
              optional: true
        - name: FIREFLY_MCP_HTTP_RATE_BURST
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: http-rate-burst
              optional: true
        # CORS
        - name: FIREFLY_MCP_HTTP_ALLOWED_ORIGINS
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: http-allowed-origins
              optional: true
        # Note: No API token env vars needed - clients pass their own tokens via Authorization header
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "256Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 5
          periodSeconds: 30
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: http
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
```

### 5. Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: firefly-mcp
  namespace: firefly-mcp
  labels:
    app: firefly-mcp
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: http
    protocol: TCP
    name: http
  selector:
    app: firefly-mcp
```

### 6. Ingress (Optional)

For external access via HTTPS:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: firefly-mcp
  namespace: firefly-mcp
  annotations:
    # For cert-manager
    cert-manager.io/cluster-issuer: letsencrypt-prod
    # For nginx ingress controller
    nginx.ingress.kubernetes.io/proxy-body-size: "10m"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - mcp.example.com
    secretName: firefly-mcp-tls
  rules:
  - host: mcp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: firefly-mcp
            port:
              name: http
```

## Quick Deploy

Create a single file `firefly-mcp.yaml` with all manifests separated by `---`, then:

```bash
# Create namespace
kubectl create namespace firefly-mcp

# Create configmap (replace with your Firefly III URL)
kubectl create configmap firefly-mcp-config \
  --namespace firefly-mcp \
  --from-literal=server-url="https://firefly.example.com/api" \
  --from-literal=log-level="info" \
  --from-literal=client-timeout="30" \
  --from-literal=http-read-timeout="30" \
  --from-literal=http-write-timeout="30" \
  --from-literal=http-idle-timeout="120" \
  --from-literal=http-session-timeout="300" \
  --from-literal=http-rate-limit="10.0" \
  --from-literal=http-rate-burst="20" \
  --from-literal=http-allowed-origins="*"

# Apply manifests
kubectl apply -f firefly-mcp.yaml

# Check deployment status
kubectl -n firefly-mcp get pods
kubectl -n firefly-mcp logs -l app=firefly-mcp
```

## Using with MCP Clients

Once deployed, configure your MCP client to connect via HTTP.

**Important**: Each client must provide their own Firefly III Personal Access Token in the `Authorization` header.

### Claude Desktop Configuration

```json
{
  "mcpServers": {
    "firefly-iii": {
      "url": "http://mcp.example.com/",
      "headers": {
        "Authorization": "Bearer <your-firefly-iii-personal-access-token>"
      }
    }
  }
}
```

Generate your Personal Access Token in Firefly III: **Profile → OAuth → Personal Access Tokens**

### Internal Cluster Access

For services within the cluster:

```
http://firefly-mcp.firefly-mcp.svc.cluster.local/
```

Remember to include the Authorization header with each client's Firefly III token.

## Health Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Liveness probe - returns 200 if server is running |
| `GET /ready` | Readiness probe - returns 200 if server can handle requests |
| `GET /` | MCP protocol endpoint (POST for tool calls) |

## Rate Limiting

The server includes built-in per-IP rate limiting:

- **Default**: 10 requests/second with burst capacity of 20
- **Per-IP tracking**: Each client IP gets independent limits
- **Proxy support**: Respects `X-Forwarded-For` and `X-Real-IP` headers
- **Health endpoints excluded**: `/health` and `/ready` bypass rate limiting

Configure via:
- `FIREFLY_MCP_HTTP_RATE_LIMIT` - requests per second (e.g., `10.0`)
- `FIREFLY_MCP_HTTP_RATE_BURST` - burst capacity (e.g., `20`)

When rate limit is exceeded, the server returns HTTP 429 (Too Many Requests).

## Scaling Considerations

- The MCP server is stateless and can be horizontally scaled
- Use `replicas: 2+` for high availability
- Built-in rate limiting protects the Firefly III API
- Consider additional rate limiting at Ingress level for DDoS protection

## Troubleshooting

### Check pod status
```bash
kubectl -n firefly-mcp get pods
kubectl -n firefly-mcp describe pod -l app=firefly-mcp
```

### View logs
```bash
kubectl -n firefly-mcp logs -l app=firefly-mcp --tail=100 -f
```

### Test connectivity
```bash
# Port forward for local testing
kubectl -n firefly-mcp port-forward svc/firefly-mcp 8080:80

# Test health endpoint
curl http://localhost:8080/health

# Test MCP endpoint (should return error without proper MCP request)
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>"
```

### Common Issues

1. **CrashLoopBackOff**: Check logs for configuration errors (missing `FIREFLY_MCP_SERVER_URL`)
2. **ImagePullBackOff**: Verify image name and registry credentials
3. **Connection refused to Firefly III**: Verify `FIREFLY_MCP_SERVER_URL` is accessible from the cluster
4. **401 Unauthorized**: Client is not providing valid `Authorization: Bearer <token>` header
5. **429 Too Many Requests**: Rate limit exceeded, adjust `FIREFLY_MCP_HTTP_RATE_LIMIT` if needed
