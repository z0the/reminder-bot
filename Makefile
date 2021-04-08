.PHONY:
.SILENT:

build:
	go build -o ./bin/bot cmd/main.go
run: build
	./.bin/bot

build-image:
	docker build -t tg-bot-reminder .
start-container:
	docker run --name tg-bot -p 80:80 --env-file .env tg-bot-reminder