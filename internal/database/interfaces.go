package database

import (
	"reminder-bot/internal/models"
)

type BotDataBase interface {
	// Reminds logic
	CreateRemind(remind models.Remind) (models.Remind, error)
	GetLastRemindIDByChatID(id int64) (int, error)
	GetAllRemindes() ([]models.Remind, error)
	GetAllRemindesByChatID(id int64) ([]models.Remind, error)
	UpdateRemind(remind models.Remind, key string, value interface{}) error

	// Users logic
	CreateUser(user models.User) error
	GetAllUsers() ([]models.User, error)
	GetUserByChatID(id int64) (models.User, error)
}
