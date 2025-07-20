# Multi-stage build for Ardilea LLM Engine
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy and build engine
COPY engine/ ./engine/
WORKDIR /app/engine
RUN CGO_ENABLED=0 GOOS=linux go build -o ardilea-engine .

# Build BASIC interpreter
WORKDIR /app
COPY basic_reference_impl.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o basic basic_reference_impl.go

# Final runtime image
FROM alpine:3.18

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /workspace

# Copy binaries
COPY --from=builder /app/engine/ardilea-engine /usr/local/bin/
COPY --from=builder /app/basic ./

# Copy workspace files
COPY test_runner.go ./
COPY tests/ ./tests/
COPY CLAUDE.md ./
COPY config.json ./

# Set permissions
RUN chmod +x /usr/local/bin/ardilea-engine ./basic

# Run the engine
CMD ["ardilea-engine"]