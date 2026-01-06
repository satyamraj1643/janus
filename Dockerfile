
# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o janus-service ./cmd/api/main.go

# Runtime Stage (Minimal image)
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/janus-service .

# Copy config files (if needed, though ideally config is in DB)
# COPY --from=builder /app/config ./config 

# Expose port
EXPOSE 8080

# Run
CMD ["./janus-service"]
