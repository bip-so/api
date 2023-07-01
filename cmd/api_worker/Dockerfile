FROM golang:1.18-alpine3.15 as builder
WORKDIR /app

COPY . .
RUN env GOOS=linux GOARCH=amd64 go build -o bip-api-worker cmd/api_worker/main.go

FROM alpine:3.15.0

WORKDIR /bip
COPY --from=builder /app/bip-api-worker ./bip-api-worker
COPY --from=builder /app/.env ./.env
EXPOSE 9003
CMD [ "/bip/bip-api-worker" ]
