package postgresdb

import (
	"reminder-bot/internal/models"

	"gorm.io/gorm"
)


type Remind struct {
	gorm.Model
	models.Remind
}
