FROM golang:1.26.2-bookworm AS builder

RUN apt-get update \
  && apt-get install -y --no-install-recommends librdkafka-dev pkg-config \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /wallet_core ./cmd/wallet_core

FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends librdkafka1 ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /wallet_core /app/wallet_core

EXPOSE 8080

CMD ["/app/wallet_core"]
