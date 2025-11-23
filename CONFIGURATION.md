# Configuration Guide

This guide provides comprehensive information about configuring the Firefly III MCP Server.

## Table of Contents

- [Overview](#overview)
- [Configuration Methods](#configuration-methods)
- [Configuration Options](#configuration-options)
- [Environment Variables](#environment-variables)
- [Configuration Precedence](#configuration-precedence)
- [Deployment Scenarios](#deployment-scenarios)
- [Validation](#validation)
- [Troubleshooting](#troubleshooting)

## Overview

The Firefly III MCP Server supports flexible configuration through:
1. **YAML configuration files** - Best for development and local testing
2. **Environment variables** - Best for production, containers, and CI/CD
3. **Hybrid approach** - Combine both methods with environment variables overriding YAML values

### Key Benefits

- **Security**: Keep sensitive data like API tokens out of version control
- **Flexibility**: Easy configuration for different environments (dev, staging, prod)
- **Container-friendly**: Perfect for Docker, Kubernetes, and cloud deployments
- **CI/CD ready**: Integrate with secret management systems

## Configuration Methods

### Method 1: YAML Configuration File

Create a `config.yaml` file in the same directory as the server binary:

```yaml
server:
  url: https://your-firefly-instance.com/api

api:
  token: your-personal-access-token

client:
  timeout: 30

limits:
  accounts: 100
  transactions: 100
  categories: 100
  budgets: 100

mcp:
  name: firefly-iii-mcp
  version: 1.0.0
  instructions: MCP server for Firefly III personal finance management
```

**Usage:**
```bash
# Use default config.yaml
./mcp-server

# Use custom config file
./mcp-server /path/to/custom-config.yaml
```

### Method 2: Environment Variables

Set environment variables with the `FIREFLY_MCP_` prefix:

```bash
export FIREFLY_MCP_SERVER_URL="https://your-firefly-instance.com/api"
export FIREFLY_MCP_API_TOKEN="your-personal-access-token"
export FIREFLY_MCP_CLIENT_TIMEOUT="60"
export FIREFLY_MCP_LIMITS_ACCOUNTS="200"
```

**Usage:**
```bash
# Environment variables will be used automatically
./mcp-server
```

### Method 3: Hybrid Configuration

Use YAML for default values and environment variables for overrides:

**config.yaml:**
```yaml
server:
  url: https://dev.firefly.local/api
api:
  token: dev-token
client:
  timeout: 30
```

**Override in production:**
```bash
export FIREFLY_MCP_SERVER_URL="https://prod.firefly.com/api"
export FIREFLY_MCP_API_TOKEN="prod-token"
export FIREFLY_MCP_CLIENT_TIMEOUT="90"
```

The server will use:
- `server.url` from environment variable (prod.firefly.com)
- `api.token` from environment variable (prod-token)
- `client.timeout` from environment variable (90)
- `limits.*` from YAML (since not overridden)
- `mcp.*` from YAML (since not overridden)

## Configuration Options

### Server Configuration

#### `server.url` (Required)

The base URL of your Firefly III API endpoint.

- **Type**: String
- **Required**: Yes
- **Environment Variable**: `FIREFLY_MCP_SERVER_URL`
- **Example**: `https://firefly.example.com/api`
- **Notes**:
  - Must include `/api` at the end
  - Must use HTTPS in production
  - Should not include trailing slash

### API Configuration

#### `api.token` (Required)

Your Firefly III Personal Access Token for API authentication.

- **Type**: String
- **Required**: Yes
- **Environment Variable**: `FIREFLY_MCP_API_TOKEN`
- **How to obtain**:
  1. Log into your Firefly III instance
  2. Go to Profile → OAuth → Personal Access Tokens
  3. Click "Create New Token"
  4. Copy the generated token
- **Security**: Never commit this to version control!

### Client Configuration

#### `client.timeout`

HTTP request timeout in seconds.

- **Type**: Integer
- **Required**: No
- **Default**: 30
- **Environment Variable**: `FIREFLY_MCP_CLIENT_TIMEOUT`
- **Range**: 1-600 (1 second to 10 minutes)
- **Recommended**:
  - Development: 30-60 seconds
  - Production: 60-120 seconds
  - Slow networks: 120+ seconds

### Limits Configuration

These settings control the maximum number of items returned per API request.

#### `limits.accounts`

Maximum number of accounts to fetch per request.

- **Type**: Integer
- **Required**: No
- **Default**: 100
- **Environment Variable**: `FIREFLY_MCP_LIMITS_ACCOUNTS`
- **Range**: 1-1000

#### `limits.transactions`

Maximum number of transactions to fetch per request.

- **Type**: Integer
- **Required**: No
- **Default**: 100
- **Environment Variable**: `FIREFLY_MCP_LIMITS_TRANSACTIONS`
- **Range**: 1-1000

#### `limits.categories`

Maximum number of categories to fetch per request.

- **Type**: Integer
- **Required**: No
- **Default**: 100
- **Environment Variable**: `FIREFLY_MCP_LIMITS_CATEGORIES`
- **Range**: 1-1000

#### `limits.budgets`

Maximum number of budgets to fetch per request.

- **Type**: Integer
- **Required**: No
- **Default**: 100
- **Environment Variable**: `FIREFLY_MCP_LIMITS_BUDGETS`
- **Range**: 1-1000

### MCP Configuration

These settings configure the MCP server metadata.

#### `mcp.name`

The name of the MCP server.

- **Type**: String
- **Required**: No
- **Default**: `firefly-iii-mcp`
- **Environment Variable**: `FIREFLY_MCP_MCP_NAME`

#### `mcp.version`

The version of the MCP server.

- **Type**: String
- **Required**: No
- **Default**: `1.0.0`
- **Environment Variable**: `FIREFLY_MCP_MCP_VERSION`

#### `mcp.instructions`

Description/instructions for the MCP server.

- **Type**: String
- **Required**: No
- **Default**: `MCP server for Firefly III personal finance management`
- **Environment Variable**: `FIREFLY_MCP_MCP_INSTRUCTIONS`

## Environment Variables

### Complete List

| Variable | YAML Path | Type | Required | Default |
|----------|-----------|------|----------|---------|
| `FIREFLY_MCP_SERVER_URL` | `server.url` | string | Yes | - |
| `FIREFLY_MCP_API_TOKEN` | `api.token` | string | Yes | - |
| `FIREFLY_MCP_CLIENT_TIMEOUT` | `client.timeout` | int | No | 30 |
| `FIREFLY_MCP_LIMITS_ACCOUNTS` | `limits.accounts` | int | No | 100 |
| `FIREFLY_MCP_LIMITS_TRANSACTIONS` | `limits.transactions` | int | No | 100 |
| `FIREFLY_MCP_LIMITS_CATEGORIES` | `limits.categories` | int | No | 100 |
| `FIREFLY_MCP_LIMITS_BUDGETS` | `limits.budgets` | int | No | 100 |
| `FIREFLY_MCP_MCP_NAME` | `mcp.name` | string | No | firefly-iii-mcp |
| `FIREFLY_MCP_MCP_VERSION` | `mcp.version` | string | No | 1.0.0 |
| `FIREFLY_MCP_MCP_INSTRUCTIONS` | `mcp.instructions` | string | No | MCP server for... |

### Naming Convention

Environment variables follow this pattern:
```
FIREFLY_MCP_<SECTION>_<KEY>
```

Where:
- `FIREFLY_MCP` is the fixed prefix
- `<SECTION>` corresponds to the YAML section (SERVER, API, CLIENT, LIMITS, MCP)
- `<KEY>` corresponds to the configuration key (URL, TOKEN, TIMEOUT, etc.)

**Examples:**
- `server.url` → `FIREFLY_MCP_SERVER_URL`
- `api.token` → `FIREFLY_MCP_API_TOKEN`
- `limits.accounts` → `FIREFLY_MCP_LIMITS_ACCOUNTS`

## Configuration Precedence

When the same configuration is specified in multiple places, the following precedence order applies (highest to lowest):

1. **Environment Variables** (highest priority)
2. **YAML Configuration File**
3. **Default Values** (lowest priority)

### Example

Given this `config.yaml`:
```yaml
server:
  url: https://dev.firefly.local/api
client:
  timeout: 30
limits:
  accounts: 50
```

And these environment variables:
```bash
FIREFLY_MCP_SERVER_URL="https://prod.firefly.com/api"
FIREFLY_MCP_CLIENT_TIMEOUT="90"
```

The effective configuration will be:
```yaml
server:
  url: https://prod.firefly.com/api  # From environment
client:
  timeout: 90                         # From environment
limits:
  accounts: 50                        # From YAML
  transactions: 100                   # From default
  categories: 100                     # From default
  budgets: 100                        # From default
```

## Deployment Scenarios

### Local Development

**Recommended**: Use `config.yaml` for convenience

1. Copy example configuration:
   ```bash
   cp config.yaml.example config.yaml
   ```

2. Edit with your credentials:
   ```yaml
   server:
     url: http://localhost:8080/api
   api:
     token: your-dev-token
   ```

3. Add to `.gitignore`:
   ```bash
   echo "config.yaml" >> .gitignore
   ```

### Docker Container

**Recommended**: Use environment variables

**docker-compose.yml:**
```yaml
version: '3.8'
services:
  firefly-mcp:
    image: firefly-mcp:latest
    environment:
      FIREFLY_MCP_SERVER_URL: ${FIREFLY_URL}
      FIREFLY_MCP_API_TOKEN: ${FIREFLY_TOKEN}
      FIREFLY_MCP_CLIENT_TIMEOUT: "60"
    stdin_open: true
    tty: true
```

**.env file:**
```bash
FIREFLY_URL=https://firefly.example.com/api
FIREFLY_TOKEN=your-secret-token
```

**Run:**
```bash
docker-compose up
```

### Kubernetes

**Recommended**: Use Secrets and ConfigMaps

**secret.yaml:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: firefly-mcp-secrets
type: Opaque
stringData:
  api-token: your-personal-access-token
```

**configmap.yaml:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: firefly-mcp-config
data:
  server-url: "https://firefly.example.com/api"
  client-timeout: "60"
```

**deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: firefly-mcp
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: firefly-mcp
        image: firefly-mcp:latest
        env:
        - name: FIREFLY_MCP_SERVER_URL
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: server-url
        - name: FIREFLY_MCP_API_TOKEN
          valueFrom:
            secretKeyRef:
              name: firefly-mcp-secrets
              key: api-token
        - name: FIREFLY_MCP_CLIENT_TIMEOUT
          valueFrom:
            configMapKeyRef:
              name: firefly-mcp-config
              key: client-timeout
```

### CI/CD Pipeline

**Recommended**: Use secret management

**GitHub Actions example:**
```yaml
name: Deploy MCP Server
on: [push]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Build
        run: go build -o mcp-server ./cmd/mcp-server

      - name: Test with config
        env:
          FIREFLY_MCP_SERVER_URL: ${{ secrets.FIREFLY_URL }}
          FIREFLY_MCP_API_TOKEN: ${{ secrets.FIREFLY_TOKEN }}
        run: go test ./...
```

## Validation

The server validates configuration on startup and will fail with clear error messages if:

### Required Fields Missing

```
Error: server.url is required (set via config file or FIREFLY_MCP_SERVER_URL)
Error: api.token is required (set via config file or FIREFLY_MCP_API_TOKEN)
```

**Solution**: Set the required configuration values.

### Invalid Values

```
Error: client.timeout must be positive
Error: limits.accounts must be positive
```

**Solution**: Ensure numeric values are positive integers.

### File Access Errors

```
Error: failed to access config file: permission denied
```

**Solution**: Check file permissions or use environment variables instead.

## Troubleshooting

### Problem: Server can't find config file

**Symptoms:**
```
Error: server.url is required
```

**Solutions:**
1. Ensure `config.yaml` exists in the current directory
2. Specify full path: `./mcp-server /full/path/to/config.yaml`
3. Use environment variables instead

### Problem: Environment variables not being read

**Symptoms:** Values from `config.yaml` are used even though environment variables are set

**Solutions:**
1. Verify environment variable names match the `FIREFLY_MCP_` prefix exactly
2. Check for typos in variable names (case-sensitive)
3. Ensure variables are exported: `export FIREFLY_MCP_SERVER_URL=...`
4. Verify in current shell: `echo $FIREFLY_MCP_SERVER_URL`

### Problem: Authentication failures

**Symptoms:**
```
Error: 401 Unauthorized
```

**Solutions:**
1. Verify API token is correct
2. Check token hasn't expired
3. Ensure token has necessary permissions
4. Test token with curl:
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" \
        https://your-instance.com/api/v1/about
   ```

### Problem: Connection timeouts

**Symptoms:**
```
Error: context deadline exceeded
```

**Solutions:**
1. Increase `client.timeout` value
2. Check network connectivity to Firefly III instance
3. Verify Firefly III server is running and accessible
4. Check firewall rules

### Debug Configuration Loading

To see which configuration source is being used, check the server logs:

```
Loading configuration from file: config.yaml
Configuration loaded successfully
```

Or:

```
Config file not found, using environment variables and defaults
Configuration loaded successfully
```

## Best Practices

1. **Never commit secrets**: Always add `config.yaml` to `.gitignore`
2. **Use environment variables in production**: More secure and flexible
3. **Set reasonable timeouts**: Balance between responsiveness and reliability
4. **Document your setup**: Keep track of which environment variables you're using
5. **Test configuration changes**: Verify server starts successfully after changes
6. **Use secret management**: For production, integrate with Vault, AWS Secrets Manager, etc.
7. **Monitor logs**: Check server logs to confirm configuration is loaded correctly
8. **Keep tokens secure**: Rotate API tokens regularly

## Additional Resources

- [Firefly III API Documentation](https://docs.firefly-iii.org/api/)
- [Viper Configuration Library](https://github.com/spf13/viper)
- [12-Factor App Config](https://12factor.net/config)
- [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
