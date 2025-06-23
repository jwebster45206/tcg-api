# Build stage
FROM golang:1.24.3-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies with current versions
RUN apk add --no-cache \
    ca-certificates=20241121-r2

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tcg-api ./cmd/tcg-api

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/tcg-api .

EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["./tcg-api"]
