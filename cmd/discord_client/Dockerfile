FROM golang:1.18-alpine3.15 as builder
WORKDIR /app

COPY . .
RUN env GOOS=linux GOARCH=amd64 go build -o bip-be-discord-client /app/cmd/discord_client/main.go

FROM alpine:3.15.0

WORKDIR /bip
COPY --from=builder /app/bip-be-discord-client ./bip-be-discord-client
COPY --from=builder /app/.env ./.env
EXPOSE 9002
CMD [ "/bip/bip-be-discord-client" ]
