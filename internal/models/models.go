package models

import (
	"time"
)

type Remind struct {
	ID             uint `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Text           string
	ActivationTime time.Time
	ChatID         int64
	IDForChat      int
	ServingNow     bool
	AlreadyServed  bool
}

type User struct {
	ID             uint `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	ChatID         int64
	TimeZoneOffset int
	AutoDelete     bool
}
