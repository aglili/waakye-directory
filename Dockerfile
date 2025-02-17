# Build Stage
FROM golang:1.23-alpine as builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final Image Stage
FROM alpine:3.18

WORKDIR /app

# Install netcat and other useful tools for debugging
RUN apk add --no-cache netcat-openbsd

# Copy the built binary and migrate tool
COPY --from=builder /app/main ./main 
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copy migrations and scripts
COPY migrations ./migrations
COPY scripts ./scripts

# Make sure entry_point.sh is executable
COPY scripts/entry_point.sh /scripts/entry_point.sh
RUN chmod +x /scripts/entry_point.sh

# Add debugging commands
RUN ls -la /app
RUN pwd

# Set entrypoint to the script
ENTRYPOINT ["/scripts/entry_point.sh"]