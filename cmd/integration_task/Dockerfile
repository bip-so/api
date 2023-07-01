FROM golang:1.18-alpine3.15 as builder
WORKDIR /app

COPY . .
RUN env GOOS=linux GOARCH=amd64 go build -o bip-integration-task cmd/integration_task/main.go

FROM alpine:3.15.0

WORKDIR /bip
COPY --from=builder /app/bip-integration-task ./bip-integration-task
COPY --from=builder /app/.env ./.env
EXPOSE 9003
CMD [ "/bip/bip-integration-task" ]
