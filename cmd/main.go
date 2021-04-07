package main

import (
	"os"
	"reminder-bot/internal/database/postgresdb"
	"reminder-bot/internal/models"
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
	db := postgresdb.NewPostgresDB(postgresdb.StartConfig{
		Host:     viper.GetString("dbCon.host"),
		Port:     viper.GetString("dbCon.port"),
		Username: viper.GetString("dbCon.username"),
		DBName:   viper.GetString("dbCon.dbname"),
		SSLMode:  viper.GetString("dbCon.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	db.AutoMigrate(&models.Remind{})
	db.AutoMigrate(&models.User{})
	botDB := postgresdb.NewBotDataBase(db)
	token := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logrus.Fatal(err)
	}
	bot.Debug = false
	telegramBot := telegram.NewBot(bot, botDB, os.Getenv("AIzaSyDH-LutBNEWdECsnCgTKoNRdbRTXdfBCw0"))
	telegramBot.Start()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
