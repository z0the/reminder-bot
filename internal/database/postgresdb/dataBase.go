package postgresdb

import (
	"reminder-bot/internal/models"

	"gorm.io/gorm"
)

type DataBase struct {
	db *gorm.DB
}

func NewBotDataBase(db *gorm.DB) *DataBase {
	return &DataBase{db: db}
}
func (p *DataBase) CreateRemind(remind models.Remind) error {
	result := p.db.Create(&remind)
	return result.Error
}
func (p *DataBase) GetRemindByID(id int) error {
	var remind Remind
	result := p.db.First(&remind, "where id=?", id)
	return result.Error
}
func (p *DataBase) GetAllRemindes() error {
	remindes := make([]Remind, 10)
	res := p.db.Find(&remindes)
	return res.Error
}
