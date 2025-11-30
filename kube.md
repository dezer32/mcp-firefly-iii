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

Non-sensitive configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: firefly-mcp-config
  namespace: firefly-mcp
data:
  # Firefly III API URL (the base URL that clients' tokens will authenticate against)
  server-url: "https://firefly.example.com/api"
  # Log level: debug, info, warn, error
  log-level: "info"
  # HTTP timeout in seconds
  client-timeout: "30"
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
        # Optional: Log level
        - name: FIREFLY_MCP_LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: log-level
              optional: true
        # Optional: HTTP client timeout
        - name: FIREFLY_MCP_CLIENT_TIMEOUT
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: client-timeout
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
  --from-literal=client-timeout="30"

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
      "url": "http://mcp.example.com/mcp",
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
http://firefly-mcp.firefly-mcp.svc.cluster.local/mcp
```

Remember to include the Authorization header with each client's Firefly III token.

## Health Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Liveness probe - returns 200 if server is running |
| `GET /ready` | Readiness probe - returns 200 if server can handle requests |

## Scaling Considerations

- The MCP server is stateless and can be horizontally scaled
- Use `replicas: 2+` for high availability
- Consider rate limiting at Ingress level to protect Firefly III API

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
```

### Common Issues

1. **CrashLoopBackOff**: Check logs for configuration errors (missing env vars)
2. **ImagePullBackOff**: Verify image name and registry credentials
3. **Connection refused to Firefly III**: Verify `FIREFLY_MCP_SERVER_URL` is accessible from the cluster
