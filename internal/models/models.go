package models

import (
	"time"

	"gorm.io/gorm"
)

type Remind struct {
	gorm.Model
	Text           string
	ActivationTime time.Time
	ChatID         int64
	IDForChat      int
	ServingNow     bool
	AlreadyServed  bool
}

type User struct{
	gorm.Model
	Name string
	ChatID int64
	TimeZoneOffset int
}