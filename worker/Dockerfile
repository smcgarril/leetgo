# Step 1: Build the Go application
ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /app .

FROM golang:1.22

WORKDIR /app

COPY --from=builder /app /app/worker
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

EXPOSE 8081

CMD ["/app/worker/leetgo-worker", "serve", "--http=0.0.0.0:8081"]