package model

import "gorm.io/gorm"

type RemindCreator interface {
}
type RemindList interface {
}
type Remind interface {
}
type Model struct {
	RemindCreator
	RemindList
	Remind
}

func NewModel(db *gorm.DB) *Model {
	return &Model{}
}
