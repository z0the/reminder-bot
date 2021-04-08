package postgresdb

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type StartConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg StartConfig, dev bool) *gorm.DB {
	var dsn string
	if dev {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	} else {
		dsn = fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Username, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color
		},
	)
	DBConfig := &gorm.Config{
		Logger: newLogger,
	}
	db, err := gorm.Open(postgres.Open(dsn), DBConfig)
	if err != nil {
		logrus.Fatal("Error on connection to db, err: ", err)
	}
	logrus.Info("Successful database connection!")
	return db
}
