package database

import (
	"reminder-bot/internal/models"
)

type BotDataBase interface {
	CreateRemind(remind models.Remind) error
	GetRemindByID(id int) error
	GetAllRemindes() error
}

// func NewModel(db *gorm.DB) *Model {
// 	return &Model{
// 		RemindCreator: NewCreatePostgres(db),
// 	}
// }
