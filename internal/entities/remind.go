package entities

import (
	"time"

	"gorm.io/gorm"
)


type Remind struct {
	gorm.Model
	Text string
	ChatID int64
	ActivationTime time.Time
}