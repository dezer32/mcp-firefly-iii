# GitHub Copilot Instructions for Firefly III MCP Server

## Project Overview

This is a Model Context Protocol (MCP) server implementation for Firefly III personal finance management system. The server provides a bridge between AI assistants (like Claude) and Firefly III instances, enabling programmatic access to financial data through a standardized protocol.

## Tech Stack

- **Language**: Go 1.24.5
- **MCP SDK**: `modelcontextprotocol/go-sdk` v0.2.0
- **HTTP Client**: Generated via `oapi-codegen`
- **Testing**: `stretchr/testify`
- **Configuration**: `gopkg.in/yaml.v3`

## Coding Guidelines

### Testing
- Write unit tests for mapper functions
- Write integration tests for MCP tools (requires live Firefly III instance)
- Use environment variables for test credentials
- Run `go test ./pkg/fireflyMCP -v` before committing

### Code Style
- Run `go fmt ./...` before committing
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and testable

### Data Transfer
- **CRITICAL**: Always use the DTOs defined in `pkg/fireflyMCP/dto.go`
- Never return raw API responses directly
- Map API responses to simplified DTOs for cleaner output
- Include pagination support in list responses

### Error Handling
- Return MCP error responses, not Go errors
- Include helpful error messages
- Log errors for debugging
- Handle nil pointers gracefully

## Issue Tracking with bd

**CRITICAL**: This project uses **bd (beads)** for ALL task tracking. Do NOT create markdown TODO lists.

### Essential Commands

```bash
# Find work
bd ready --json                    # Unblocked issues

# Create and manage
bd create "Title" -t bug|feature|task -p 0-4 --json
bd update <id> --status in_progress --json
bd close <id> --reason "Done" --json

# Search
bd list --status open --priority 1 --json
bd show <id> --json
```

### Workflow

1. **Check ready work**: `bd ready --json`
2. **Claim task**: `bd update <id> --status in_progress`
3. **Work on it**: Implement, test, document
4. **Discover new work?** `bd create "Found bug" -p 1 --deps discovered-from:<parent-id> --json`
5. **Complete**: `bd close <id> --reason "Done" --json`
6. **Commit together**: Always commit the `.beads/issues.jsonl` file together with the code changes

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (default, nice-to-have)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

## Project Structure

```
firefly-iii/
├── cmd/mcp-server/         # Main entry point
│   └── main.go
├── pkg/
│   ├── client/            # Auto-generated API client
│   │   ├── client.go
│   │   └── generate.go
│   └── fireflyMCP/        # MCP server implementation
│       ├── config.go      # Configuration management
│       ├── dto.go         # Data transfer objects
│       ├── server.go      # MCP server & handlers
│       ├── integration_test.go  # Integration tests
│       └── mapper_test.go       # Unit tests for mappers
├── resources/             # OpenAPI specifications
├── example/              # Example usage
├── config.yaml          # Configuration file
├── README.md           # User documentation
├── TESTING.md          # Testing documentation
├── CLAUDE.md          # Development guide
└── AGENTS.md          # AI agent instructions
```

## Development Commands

### Building
```bash
go build -o mcp-server ./cmd/mcp-server
```

### Running
```bash
./mcp-server                          # Use default config.yaml
./mcp-server /path/to/custom-config.yaml
```

### Code Generation
```bash
cd pkg/client
go generate  # Regenerate API client from OpenAPI spec
```

### Testing
```bash
go test ./...                         # All tests
go test -cover ./pkg/fireflyMCP      # With coverage
go test -v -run TestMapBudget ./pkg/fireflyMCP  # Specific test
```

## Important Rules

- ✅ Use bd for ALL task tracking
- ✅ Always use DTOs from `dto.go` for data responses
- ✅ Write tests for new functionality
- ✅ Use `--json` flag for programmatic bd commands
- ❌ Do NOT create markdown TODO lists
- ❌ Do NOT return raw API responses
- ❌ Do NOT commit sensitive tokens

---

**For detailed workflows and advanced features, see [AGENTS.md](../AGENTS.md) and [CLAUDE.md](../CLAUDE.md)**
