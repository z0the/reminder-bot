FROM golang:1.16.3-alpine3.12 AS builder

COPY . /reminder-bot/
WORKDIR /reminder-bot/

RUN go mod download
RUN go build -o ./bin/bot cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /reminder-bot/bin/bot .
COPY --from=0 /reminder-bot/configs configs/

CMD ["./bot"]
