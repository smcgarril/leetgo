ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /app .

FROM golang:1.22

# Comment out for local development #
# ENV WORKER_HOST=https://leetgo-worker.fly.dev
# ENV WORKER_PORT=443
# ENV WORKER_PATH=/process-code

WORKDIR /app

COPY --from=builder /app /app/server
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

EXPOSE 8080

COPY ./public /app/public
COPY ./db /app/db
CMD ["/app/server/leetgo", "serve", "--http=0.0.0.0:8080"]