.PHONY:
.SILENT:

build:
	go build -o ./bin/bot cmd/main.go
run: build
	./.bin/bot

build-image:
	docker build -t zothe/telegram-bot:latest.
start-container:
	docker run --name tg-bot -p 80:80 --env-file .env zothe/telegram-bot:latest