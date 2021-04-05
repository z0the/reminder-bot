package main

import (
	"os"
	"reminder-bot/internal/database/postgresdb"
	"reminder-bot/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, DisableTimestamp: false, FullTimestamp: true})
	logrus.Info("Initializing App...")
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
	err = initConfig()
	if err != nil {
		logrus.Fatal("Error loading configs file")
	}
}
func main() {
	db := postgresdb.NewPostgresDB(postgresdb.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	db.AutoMigrate(&postgresdb.Remind{})
	botDB := postgresdb.NewBotDataBase(db)
	token := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logrus.Fatal(err)
	}
	bot.Debug = false
	telegramBot := telegram.NewBot(bot, botDB)
	telegramBot.Start()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
