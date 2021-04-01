package dao

import (
	"read-test-server/common/gorm2/mysql"
	"read-test-server/model"
)

var db mysql.DB

func InitMysqlDb(dbConfig *mysql.DbConfig, log mysql.CustomerLog) {
	db = mysql.MakeDB(dbConfig, log)
}

func CreateUser(user *model.User) error {
	return db.GetDB().Create(user).Error
}
