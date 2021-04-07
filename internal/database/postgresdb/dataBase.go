package postgresdb

import (
	"reminder-bot/internal/models"

	"github.com/sirupsen/logrus"
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
func (p *DataBase) UpdateRemind(remind models.Remind, key string, value interface{}) error {
	logrus.Infof("Updating remind key: %s, value: %s", key, value)
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
