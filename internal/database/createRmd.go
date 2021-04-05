package database

import (
	"reminder-bot/internal/entities"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CreatePostgres struct {
	db *gorm.DB
}

func NewCreatePostgres(db *gorm.DB) *CreatePostgres {
	return &CreatePostgres{db: db}
}
func (m *CreatePostgres) CreateRemind(remind entities.Remind) (int, error) {
	m.db.Create(&remind)
	logrus.Info("Remind created! ", remind.ID)
	return int(remind.ID), nil
}
