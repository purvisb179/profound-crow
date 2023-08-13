# Stage 1: build the Go binary
FROM golang:1.18 as builder

WORKDIR /app

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp main.go

# Stage 2: copy the Go binary to an Alpine container and update certs
FROM alpine:latest

# Update certificates in Alpine
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/myapp /app/myapp

# The command to start the app
CMD ["/app/myapp"]
