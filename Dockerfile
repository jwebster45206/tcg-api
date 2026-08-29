# Build stage
FROM golang:1.26.7-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o tcg-api ./cmd/tcg-api

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /app
COPY --from=builder /app/tcg-api .

EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["./tcg-api"]
