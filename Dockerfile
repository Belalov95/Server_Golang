FROM golang:1.23-alpine AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/config
COPY internal/config/config.yml /app/config/config.yml
COPY .env /app/.env

RUN go build -o /app/api-gateway ./cmd/
RUN ls -l /app

FROM alpine:latest 

WORKDIR /app
COPY --from=builder /app/api-gateway .
COPY --from=builder /app/config ./config
COPY --from=builder /app/.env ./.env

ENTRYPOINT ["./api-gateway"]
