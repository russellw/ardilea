# Multi-stage build for Ardilea LLM Engine
FROM golang:1.21-alpine AS builder

# Install git for potential go mod requirements
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy engine source code
COPY engine/ ./engine/

# Copy workspace files (BASIC interpreter and tests)
COPY *.go ./
COPY *.bat ./
COPY tests/ ./tests/
COPY README_testing.md ./
COPY CLAUDE.md ./

# Build the engine
WORKDIR /app/engine
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o ardilea-engine .

# Build the BASIC interpreter for the workspace
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o basic basic_reference_impl.go

# Final runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /workspace

# Copy the built engine
COPY --from=builder /app/engine/ardilea-engine /usr/local/bin/

# Copy the BASIC interpreter and test infrastructure
COPY --from=builder /app/basic ./
COPY --from=builder /app/test_runner.go ./
COPY --from=builder /app/tests/ ./tests/
COPY --from=builder /app/README_testing.md ./
COPY --from=builder /app/CLAUDE.md ./

# Create config directory
RUN mkdir -p /workspace/config

# Set executable permissions
RUN chmod +x /usr/local/bin/ardilea-engine
RUN chmod +x ./basic

# Create a default config file
RUN echo '{\n  "ollama_server": "192.168.0.63:11434",\n  "model_name": "qwen2.5:32b",\n  "workspace_dir": "/workspace"\n}' > config.json

# Run the engine
CMD ["ardilea-engine"]