FROM golang:1.22
WORKDIR /app

COPY go.* /app
RUN go mod download && go mod verify

COPY . /app
RUN go build -v -o main .

FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server

CMD ["/app/main"]