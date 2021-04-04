package main

import (
	"log"
	"os"
	"reminder-bot/internal/handler"
	"reminder-bot/internal/model"
	"reminder-bot/internal/service"
	"reminder-bot/internal/service/keyboard"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	db := model.NewPostgresDB(model.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "p@ssw0rd",
		DBName:   "botdb",
		SSLMode:  "disable",
	})
	model := model.NewModel(db)
	service := service.NewService(model)
	handler := handler.NewHandler(service)
	token := os.Getenv("BOT_TOKEN")
	bot, err := botapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := botapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("Get updates error: ", err)
	}
	keyboard := keyboard.GetMainKeyboard()
	r := handler.InitRoutes(bot, keyboard, &updates)
	for update := range updates {
		if update.Message != nil {
			r.Handle(&update)
		}
	}
}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
