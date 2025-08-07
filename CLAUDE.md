# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Mods is a command-line tool that adds AI capabilities to command pipelines. It works by reading standard input and prefacing it with a prompt supplied in the arguments, sending the text to an LLM, and printing the response. The tool supports multiple AI providers including OpenAI, Anthropic, Cohere, Google, Groq, Ollama, and local models.

## Key Architecture Components

### Core Files
- `main.go` - CLI entry point and cobra command setup
- `mods.go` - Main application logic and state machine (startState, configLoadedState, requestState, responseState, doneState, errorState)
- `config.go` - Configuration management using YAML with embedded template
- `db.go` - SQLite database operations for conversation storage
- `mcp.go` - Model Context Protocol (MCP) server integration
- `messages.go` - Message formatting and processing
- `stream.go` - Streaming response handling

### Internal Packages
- `internal/anthropic/` - Anthropic API integration
- `internal/openai/` - OpenAI API integration  
- `internal/cohere/` - Cohere API integration
- `internal/google/` - Google Gemini API integration
- `internal/ollama/` - Ollama local model integration
- `internal/copilot/` - GitHub Copilot integration
- `internal/cache/` - Conversation caching and management
- `internal/proto/` - Protocol buffer definitions
- `internal/stream/` - Streaming utilities

## Common Development Commands

### Build and Test
```bash
# Build the binary
go build -o mods

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/cache

# Run tests with verbose output
go test -v ./...

# Build with race detection (for development)
go build -race -o mods
```

### Development Setup
```bash
# Install dependencies
go mod download

# Verify dependencies
go mod verify

# Clean module cache
go clean -modcache

# Run with development flags
go run . --help
```

## Configuration System

The application uses a YAML configuration file with extensive model definitions. Configuration is managed through:

- `config_template.yml` - Embedded template with all available options
- User config file located in XDG config directory
- Environment variables for API keys (OPENAI_API_KEY, ANTHROPIC_API_KEY, etc.)

Key configuration sections:
- `apis` - API provider configurations with base URLs and models
- `roles` - Role-based system prompts and MCP tool filtering (see Role System below)
- `mcp-servers` - Model Context Protocol server configurations
- Model aliases and fallbacks for graceful degradation

## Database Schema

Uses SQLite for local conversation storage with tables:
- `conversations` - Stores conversation metadata (id, title, updated_at)
- Messages are stored as JSON within conversations

## Key Features

- **Multi-provider support** - Works with OpenAI, Anthropic, Cohere, Google, Groq, Ollama, and local models
- **Conversation persistence** - Local SQLite database with SHA-1 conversation IDs
- **MCP integration** - Model Context Protocol for extending model capabilities
- **Streaming responses** - Real-time response display with Bubble Tea TUI
- **Role-based MCP filtering** - Fine-grained control over MCP tools and servers per role
- **Pipeline integration** - Reads from stdin for command line workflows

## Testing Guidelines

- Unit tests are located alongside source files with `_test.go` suffix
- Use `testdata/` directories for test fixtures (e.g., `internal/proto/testdata/`)
- Golden file testing pattern used in some packages
- Mock external API calls in tests

## Role System

### Role Configuration Structure

Roles support both backwards-compatible simple format and new structured format:

**Simple Format (legacy):**
```yaml
roles:
  shell: ["you are a shell expert", "output only commands"]
```

**Structured Format (new):**
```yaml
roles:
  shell:
    prompt:
      - "you are a shell expert"
      - "output only commands"
    allowed_servers: ["filesystem", "shell"]
    blocked_servers: ["network"] 
    allowed_tools: ["filesystem_*", "shell_*"]
    blocked_tools: ["dangerous_*", "*_delete"]
```

### MCP Tool Filtering

- **Tool Naming Convention**: `{server_name}_{tool_name}` (e.g., `github_list_repos`)
- **Glob Pattern Support**: Use `*` wildcards in allowed/blocked lists
- **Precedence**: Blocked lists always override allowed lists
- **Filtering Levels**:
  - Server-level: `allowed_servers` and `blocked_servers`
  - Tool-level: `allowed_tools` and `blocked_tools` (matches full tool name or tool name alone)

### Key Functions

- `mcpToolsForRole(ctx, role)` - Gets filtered tools for specific role
- `isServerAllowedForRole(serverName, roleConfig)` - Server filtering logic
- `isToolAllowedForRole(toolName, serverName, roleConfig)` - Tool filtering logic
- `RoleConfig.UnmarshalYAML()` - Backwards compatibility handler

## Code Style Notes

- Uses Go 1.24 features and modules
- Follows standard Go project layout with `internal/` packages
- Extensive use of Cobra for CLI management
- Bubble Tea for terminal UI components
- Error handling with wrapped errors and specific error types