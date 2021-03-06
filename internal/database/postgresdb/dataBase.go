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
func (p *DataBase) CreateRemind(remind models.Remind) (models.Remind, error) {
	result := p.db.Create(&remind)
	return remind, result.Error
}
func (p *DataBase) DeleteRemind(id int) error {
	var remind models.Remind
	result := p.db.Delete(&remind, id)
	return result.Error
}
func (p *DataBase) GetLastRemindIDByChatID(id int64) (int, error) {
	var remind models.Remind
	queryResult := p.db.Where("chat_id=?", id).Last(&remind)
	if queryResult.Error == gorm.ErrRecordNotFound {
		return remind.IDForChat, nil
	} else {
		return remind.IDForChat, queryResult.Error
	}

}
func (p *DataBase) GetAllRemindes() ([]models.Remind, error) {
	reminds := make([]models.Remind, 10)
	queryResult := p.db.Find(&reminds)
	return reminds, queryResult.Error
}
func (p *DataBase) GetAllRemindesByChatID(id int64) ([]models.Remind, error) {
	reminds := make([]models.Remind, 10)
	queryResult := p.db.Find(&reminds, "chat_id=?", id)
	return reminds, queryResult.Error
}
func (p *DataBase) UpdateRemind(remind models.Remind, key string, value interface{}) error {
	queryResult := p.db.Model(&remind).Update(key, value)
	return queryResult.Error
}
func (p *DataBase) CreateUser(user models.User) error {
	queryResult := p.db.Create(&user)
	return queryResult.Error
}
func (p *DataBase) GetAllUsers() ([]models.User, error) {
	users := make([]models.User, 10)
	queryResult := p.db.Find(&users)
	return users, queryResult.Error
}
func (p *DataBase) GetUserByChatID(id int64) (models.User, error) {
	var user models.User
	queryResult := p.db.First(&user, "chat_id=?", id)
	return user, queryResult.Error
}
func (p *DataBase) UpdateUser(user models.User) error {
	queryResult := p.db.Model(&user).Updates(user)
	return queryResult.Error
}
func (p *DataBase) UpdateUserBool(user models.User, key string, value bool) error {
	queryResult := p.db.Model(&user).Update(key, value)
	return queryResult.Error
}
