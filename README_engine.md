# Ardilea LLM Engine

An LLM-powered agent engine that runs in a Docker container to develop and improve BASIC interpreter implementations using Ollama.

## Overview

The Ardilea Engine connects to an Ollama server to interact with large language models for automated software development. It's designed to analyze, improve, and develop BASIC interpreter code within a containerized workspace.

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Ollama server running with a compatible model
- Network access to the Ollama server

### Basic Usage

1. **Configure the engine** (optional):
   ```bash
   # Edit config.json to match your setup
   {
     "ollama_server": "192.168.0.63:11434",
     "model_name": "qwen3:30b",
     "workspace_dir": "/workspace"
   }
   ```

2. **Run with Docker Compose**:
   ```bash
   # Start the complete stack (includes Ollama server)
   docker-compose up -d

   # Or run just the engine (assumes external Ollama server)
   docker-compose up ardilea-engine
   ```

3. **Build and run manually**:
   ```bash
   # Build the container
   docker build -t ardilea-engine .

   # Run the engine
   docker run -it --rm \
     -v $(pwd)/workspace:/workspace \
     -v $(pwd)/config.json:/workspace/config.json:ro \
     ardilea-engine
   ```

## Configuration

The engine reads configuration from `config.json`:

| Setting | Default | Description |
|---------|---------|-------------|
| `ollama_server` | `192.168.0.63:11434` | Ollama server address and port |
| `model_name` | `qwen3:30b` | LLM model to use for code generation |
| `workspace_dir` | `/workspace` | Working directory inside container |

### Environment Variables

You can also use environment variables to override config:

- `OLLAMA_SERVER` - Ollama server address
- `MODEL_NAME` - Model name to use

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Docker         │    │   Ardilea       │    │    Ollama       │
│  Container      │◄──►│   Engine        │◄──►│    Server       │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │ Workspace   │ │    │ │ Go Engine   │ │    │ │ LLM Model   │ │
│ │ - BASIC     │ │    │ │ - Config    │ │    │ │ (qwen3:30b) │ │
│ │ - Tests     │ │    │ │ - API Client│ │    │ │             │ │
│ │ - Output    │ │    │ │ - Agent     │ │    │ │             │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Engine Capabilities

### Analysis Mode
- Scans existing workspace for BASIC interpreter implementation
- Analyzes code quality and identifies improvement opportunities
- Suggests specific next steps for development

### Development Mode
- Generates BASIC interpreter code from scratch
- Implements features based on test requirements
- Iterates on improvements based on test results

### Workspace Tracking
- Creates before/after snapshots of all workspace files
- Tracks added, removed, and modified files using MD5 hashing
- Generates detailed reports: `workspace-report.json` and `workspace-summary.txt`
- Displays change summary in console after completion

### Features
- **Config-driven**: Easy setup with JSON configuration
- **Ollama Integration**: Native support for Ollama API
- **Workspace Management**: Persistent workspace with volume mounting
- **Health Checking**: Verifies Ollama connectivity before starting
- **Error Handling**: Robust error handling and logging

## Workspace Structure

```
/workspace/
├── basic                    # BASIC interpreter executable
├── test_runner.go          # Test runner source
├── tests/                  # Test cases
│   ├── basic/*.bas        # BASIC test programs
│   ├── expected/*.txt     # Expected outputs
│   └── errors/*.bas       # Error test cases
├── config.json            # Engine configuration
├── workspace-report.json  # Detailed change report (generated)
└── workspace-summary.txt  # Human-readable summary (generated)
```

## Development Workflow

1. **Initial Analysis**: Engine analyzes existing codebase
2. **Gap Identification**: LLM identifies missing features or bugs
3. **Code Generation**: LLM generates improvements or new code
4. **Testing**: Engine runs test suite to verify changes
5. **Iteration**: Process repeats until goals are met

## Viewing Results

After the engine completes, check the `workspace/` directory for:

### Generated Reports
```bash
# View human-readable summary
cat workspace/workspace-summary.txt

# View detailed JSON report
cat workspace/workspace-report.json | jq .

# List all changes
grep -E "^\s*[+~-]" workspace/workspace-summary.txt
```

### Example Output
```
Workspace changes from 2024-01-20 15:30:22 to 2024-01-20 15:45:33:
- Files added: 3
- Files removed: 0  
- Files modified: 2

Added files:
  + src/new_interpreter.go
  + tests/new_test.bas
  + docs/improvements.md

Modified files:
  ~ basic_reference_impl.go (size: 15234->16789 bytes)
  ~ README.md (size: 2451->3102 bytes)
```

## Network Requirements

The engine needs network access to:
- Ollama server (default: `192.168.0.63:11434`)
- Internet (for Go module downloads during build)

## Troubleshooting

### Common Issues

1. **Connection to Ollama fails**:
   ```bash
   # Check Ollama server is running
   curl http://192.168.0.63:11434/api/tags
   
   # Verify network connectivity from container
   docker run --rm ardilea-engine ping 192.168.0.63
   ```

2. **Model not available**:
   ```bash
   # List available models on Ollama server
   curl http://192.168.0.63:11434/api/tags
   
   # Pull required model
   docker exec ollama-server ollama pull qwen3:30b
   ```

3. **Permission issues**:
   ```bash
   # Ensure workspace directory is writable
   chmod 755 workspace/
   ```

### Logs

View engine logs:
```bash
docker-compose logs -f ardilea-engine
```

## Advanced Usage

### Custom Model

To use a different model:
1. Update `config.json` with the new model name
2. Ensure the model is available on your Ollama server
3. Restart the engine

### Development Mode

For development and debugging:
```bash
# Run with interactive shell
docker run -it --rm \
  -v $(pwd):/workspace/output \
  --entrypoint /bin/sh \
  ardilea-engine
```

### Integration with CI/CD

The engine can be integrated into CI/CD pipelines for automated code improvement:

```yaml
# Example GitHub Actions workflow
- name: Run Ardilea Engine
  run: |
    docker run --rm \
      -v ${{ github.workspace }}:/workspace/output \
      ardilea-engine
```

## Security Considerations

- The engine runs as non-root user inside container
- Workspace is isolated within container
- Network access limited to Ollama server
- No persistent credentials stored in container
- Configuration mounted read-only when possible

## Contributing

1. Modify engine code in `engine/` directory
2. Test changes with Docker build
3. Update documentation as needed
4. Submit pull request

## License

This project is part of the Ardilea BASIC interpreter development toolkit.