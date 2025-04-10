# Build stage
FROM golang:1.24.2 AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the source code
COPY . ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o chatheon ./cmd/chatheon

# Final stage
FROM alpine:3.21

WORKDIR /root/

# Copy the binary from the build stage
COPY --from=builder /app/chatheon .

# Use minimal image for execution
EXPOSE 8080

CMD ["./chatheon"]