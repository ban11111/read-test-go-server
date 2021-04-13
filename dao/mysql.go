package dao

import (
	"read-test-server/common/gorm2/mysql"
	"read-test-server/model"
)

var db mysql.DB

func InitMysqlDb(dbConfig *mysql.DbConfig, log mysql.CustomerLog) {
	utilDb := mysql.MakeDBUtil(dbConfig, log)
	utilDb.CreateDB()
	db = mysql.MakeDB(dbConfig, log)
}

func QueryGlobalSettings() ([]*model.GlobalSetting, error) {
	var settings []*model.GlobalSetting
	return settings, db.GetDB().Model(&model.GlobalSetting{}).Find(&settings).Error
}

func UpdateGlobalSetting(set *model.GlobalSetting) error {
	return db.GetDB().Model(&model.GlobalSetting{}).Where("`key`=?", set.Key).UpdateColumn("value", set.Value).Error
}

func CreateUser(user *model.User) error {
	tx := db.GetDB().Begin()
	defer tx.Rollback()
	var n int64
	if err := tx.Model(&model.User{}).Where("email=?", user.Email).Count(&n).Error; err != nil {
		return err
	}
	if err := db.GetDB().Create(user).Error; err != nil {
		return err
	}
	return tx.Commit().Error
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
	return db.GetDB().Where("id=?", paper.Id).Updates(paper).Error
}

func CreatePaperSnapshot(paper *model.PaperSnapshot) error {
	return db.GetDB().Create(paper).Error
}

func QueryPaperById(id uint) (*model.Paper, error) {
	var paper model.Paper
	return &paper, db.GetDB().Model(&model.Paper{}).Where("id=?", id).Last(&paper).Error
}

func QueryPapers() ([]*model.Paper, error) {
	var papers []*model.Paper
	return papers, db.GetDB().Model(&model.Paper{}).Order("id desc").Find(&papers).Error
}

func QueryCurrentPaper() (*model.Paper, error) {
	var paper model.Paper
	return &paper, db.GetDB().Model(&model.Paper{}).Where("inuse=?", true).Last(&paper).Error
}

// 查询某个用户做了哪些试卷
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

func DeleteAnswersByUidAndPaper(uid uint, paperId uint) error {
	return db.GetDB().Model(&model.Answer{}).Where("uid=? and paper_id=?", uid, paperId).Unscoped().Delete(&model.Answer{}).Error
}

func QueryAnswerProgress(uid uint, paperId uint) (*model.Answer, error) {
	var answer model.Answer
	return &answer, db.GetDB().Model(&model.Answer{}).Where("uid=? and paper_id=?", uid, paperId).Order("word_index desc").First(&answer).Error
}
