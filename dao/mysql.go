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

func QueryUsers() ([]*model.User, error) {
	var users []*model.User
	return users, db.GetDB().Model(&model.User{}).Find(&users).Error
}

func QueryUserByEmail(email string) (*model.User, error) {
	var user model.User
	return &user, db.GetDB().Model(&user).Where("email=?", email).Last(&user).Error
}

func CreatePaper(paper *model.Paper) error {
	return db.GetDB().Create(paper).Error
}

func UpdatePaper(paper *model.Paper) error {
	return db.GetDB().Updates(paper).Error
}

func QueryPapers() ([]*model.Paper, error) {
	var papers []*model.Paper
	return papers, db.GetDB().Model(&model.Paper{}).Find(&papers).Error
}

func QueryPapersByUid(uid uint) ([]*model.Paper, error) {
	var papers []*model.Paper
	subQuery := db.GetDB().Model(&model.Answer{}).Select("distinct(paper_id)").Where("uid=?", uid)
	return papers, db.GetDB().Model(&model.Paper{}).Where("id in (?)", subQuery).Find(&papers).Error
}

func CreateAnswer(answer *model.Answer) error {
	return db.GetDB().Create(answer).Error
}

func UpdateAnswer(answer *model.Answer) error {
	return db.GetDB().Updates(answer).Error
}

func QueryAnswersByUidAndPaper(uid uint, paperId uint) ([]*model.Answer, error) {
	var answers []*model.Answer
	return answers, db.GetDB().Model(&model.Answer{}).Where("uid=? and paper_id=?", uid, paperId).Order("word_index").Find(&answers).Error
}
