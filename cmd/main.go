package main

import (
	"os"
	"reminder-bot/internal/database/postgresdb"
	"reminder-bot/internal/models"
	"reminder-bot/internal/telegram"
	"strconv"

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
		logrus.Warn(err)
	}
	err = initConfig()
	if err != nil {
		logrus.Fatal(err)
	}
}
func main() {
	devEnv, err := strconv.ParseBool(os.Getenv("DEV"))
	if err != nil {
		logrus.Fatal("Error loading DEV env, err: ", err)
	}
	db := postgresdb.NewPostgresDB(postgresdb.StartConfig{
		Host:     viper.GetString("dbCon.host"),
		Port:     viper.GetString("dbCon.port"),
		Username: viper.GetString("dbCon.username"),
		DBName:   viper.GetString("dbCon.dbname"),
		SSLMode:  viper.GetString("dbCon.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	}, devEnv)

	db.AutoMigrate(&models.Remind{})
	db.AutoMigrate(&models.User{})
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
