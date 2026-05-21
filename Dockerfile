# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

# Stage 2: Minimal run-time image
FROM alpine:latest

# Add certificates and timezone database
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /app/main .

# Copy swagger docs if they exist
COPY --from=builder /app/docs ./docs

# Create uploads folder
RUN mkdir -p uploads

EXPOSE 8080

CMD ["./main"]
