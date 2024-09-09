FROM golang:1.22.4-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o diasync ./cmd/main.go

FROM alpine

WORKDIR  /app

COPY --from=builder /build/diasync .
COPY --from=builder /build/config/config.json .

EXPOSE 8080

CMD ["./diasync", "-p", "/app/config.json"]