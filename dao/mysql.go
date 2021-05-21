package dao

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"read-test-server/common"
	"read-test-server/common/gorm2/mysql"
	"read-test-server/model"
	"regexp"
	"time"
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

func QueryUsers() ([]*model.UserInfo, error) {
	var users []*model.UserInfo
	if err := db.GetDB().Model(&model.UserInfo{}).Preload("Papers").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(uid uint) error {
	var user model.User
	if err := db.GetDB().Model(&model.User{}).Where("id=?", uid).Last(&user).Error; err != nil {
		return err
	}
	tx := db.GetDB().Begin()
	defer tx.Rollback()
	if err := tx.Model(&model.User{}).Where("id=?", uid).Unscoped().Delete(&model.User{}).Error; err != nil {
		return err
	}
	if err := tx.Model(&model.Answer{}).Where("uid=?", uid).Unscoped().Delete(&model.Answer{}).Error; err != nil {
		return err
	}
	if err := tx.Model(&model.PaperUser{}).Where("uid=?", uid).Delete(&model.PaperUser{}).Error; err != nil {
		return err
	}
	return tx.Commit().Error
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

func PublishPaper(pid uint) error {
	var count int64
	if err := db.GetDB().Model(&model.Paper{}).Where("id=?", pid).Count(&count).Error; err != nil {
		return err
	}
	if count <= 0 {
		return gorm.ErrRecordNotFound
	}
	tx := db.GetDB().Begin()
	defer tx.Rollback()
	if err := tx.Model(&model.Paper{}).Where("id <> ?", pid).UpdateColumn("inuse", false).Error; err != nil {
		return err
	}
	if err := tx.Model(&model.Paper{}).Where("id = ?", pid).UpdateColumn("inuse", true).Error; err != nil {
		return err
	}
	return tx.Commit().Error
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

func CreateAnswer(answer *model.Answer, paperName string) error {
	tx := db.GetDB().Begin()
	defer tx.Rollback()

	if err := tx.Create(answer).Error; err != nil {
		return err
	}
	if err := tx.Model(&model.PaperUser{}).Where("uid=? and pid=?", answer.Uid, answer.PaperId).FirstOrCreate(&model.PaperUser{
		Uid:   answer.Uid,
		Pid:   answer.PaperId,
		PName: paperName,
	}).Error; err != nil {
		common.Log.Error("????", zap.Error(err))
		return err
	}
	return tx.Commit().Error
}

func UpdateAnswer(answer *model.Answer) error {
	return db.GetDB().Updates(answer).Error
}

func QueryAnswersByUidAndPaper(uid uint, paperId uint) ([]*model.Answer, error) {
	var answers []*model.Answer
	return answers, db.GetDB().Model(&model.Answer{}).Where("uid=? and paper_id=?", uid, paperId).Order("word_index").Find(&answers).Error
}

func DeleteAnswersByUid(uid uint) error {
	tx := db.GetDB().Begin()
	defer tx.Rollback()

	if err := tx.Model(&model.Answer{}).Where("uid=?", uid).Unscoped().Delete(&model.Answer{}).Error; err != nil {
		return err
	}
	if err := tx.Model(&model.PaperUser{}).Where("uid=?", uid).Delete(&model.PaperUser{}).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

func DeleteAnswersByUidAndPaper(uid, paperId uint) error {
	tx := db.GetDB().Begin()
	defer tx.Rollback()

	if err := db.GetDB().Model(&model.Answer{}).Where("uid=? and paper_id=?", uid, paperId).Unscoped().Delete(&model.Answer{}).Error; err != nil {
		return err
	}
	if err := tx.Model(&model.PaperUser{}).Where("uid=? and pid=?", uid, paperId).Delete(&model.PaperUser{}).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

func QueryAnswerProgress(uid uint, paperId uint) (*model.Answer, error) {
	var answer model.Answer
	return &answer, db.GetDB().Model(&model.Answer{}).Where("uid=? and paper_id=?", uid, paperId).Order("word_index desc").First(&answer).Error
}

// For Exporter
func QueryUsersByIds(ids []uint) ([]*model.User, error) {
	var users []*model.User
	sql := db.GetDB().Model(&model.User{})
	if len(ids) > 0 {
		sql = sql.Where("id in (?)", ids)
	}
	if err := sql.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func QueryAnswersByUidsAndPaperId(uid []uint, paperId uint) ([]*model.Answer, error) {
	var answers []*model.Answer
	sql := db.GetDB().Model(&model.Answer{})
	if len(uid) > 0 {
		sql = sql.Where("uid in (?)", uid)
	}
	return answers, sql.Where("paper_id=?", paperId).Order("uid, paper_id, word_index").Find(&answers).Error
}

type countTool struct {
	Count      int
	CreateTime string
	Device     string
}

func QueryStatistics() (*model.StatisticsResp, error) {
	resp := &model.StatisticsResp{
		Chart: model.ChartInfo{
			Daily: model.ChartCount{
				NewUser:   make([]int, 7),
				NewAnswer: make([]int, 7),
				Labels:    make([]string, 7),
			},
			Monthly: model.ChartCount{
				NewUser:   make([]int, 7),
				NewAnswer: make([]int, 7),
				Labels:    make([]string, 7),
			},
		},
	}
	var count int64
	if err := db.GetDB().Model(&model.User{}).Count(&count).Error; err != nil {
		return nil, err
	}
	resp.TotalUsers = int(count)
	count = 0

	var counts []*countTool
	// =================================================================================================================
	//                                               device info
	// =================================================================================================================
	if err := db.GetDB().Model(&model.Answer{}).Select("count(*) as count, device").
		Group("device").Scan(&counts).Error; err != nil {
		return nil, err
	}
	for _, count := range counts {
		switch count.Device {
		case "desktop":
			resp.Device.Desktop = count.Count
		case "tablet":
			resp.Device.Tablet = count.Count
		case "mobile":
			resp.Device.Mobile = count.Count
		default:
			resp.Device.Unknown += count.Count
		}
	}
	resp.TotalAnswers = resp.Device.Desktop + resp.Device.Tablet + resp.Device.Mobile + resp.Device.Unknown

	// =================================================================================================================
	//                                               chart info
	// =================================================================================================================
	year, month, day := time.Now().Date()
	dailyBegin := time.Date(year, month, day, 0, 0, 0, 0, time.Local).Add(-time.Hour * 24 * 6)
	monthlyBegin := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).AddDate(0, -6, 0)

	if err := db.GetDB().Model(&model.User{}).Select("count(*) as count, @createTime as create_time").
		Where("?<=created_at", dailyBegin).Group("@createTime := date(created_at)").Scan(&counts).Error; err != nil {
		return nil, err
	}
	if err := calculateChartCounts(dailyBegin, 0, 1, counts, resp.Chart.Daily.NewUser, resp.Chart.Daily.Labels); err != nil {
		return nil, err
	}

	if err := db.GetDB().Model(&model.User{}).Select("count(*) as count, concat(@createTime, '-01') as create_time").
		Where("?<=created_at", monthlyBegin).Group("@createTime := DATE_FORMAT(created_at, '%Y-%m')").Scan(&counts).Error; err != nil {
		return nil, err
	}
	if err := calculateChartCounts(monthlyBegin, 1, 0, counts, resp.Chart.Monthly.NewUser, resp.Chart.Monthly.Labels); err != nil {
		return nil, err
	}
	resp.CurrentMonthUsers = resp.Chart.Monthly.NewUser[len(resp.Chart.Monthly.NewUser)-1]

	if err := db.GetDB().Model(&model.Answer{}).Select("count(*) as count, @createTime as create_time").
		Where("?<=created_at", dailyBegin).Group("@createTime := date(created_at)").Scan(&counts).Error; err != nil {
		return nil, err
	}
	if err := calculateChartCounts(dailyBegin, 0, 1, counts, resp.Chart.Daily.NewAnswer, nil); err != nil {
		return nil, err
	}

	if err := db.GetDB().Model(&model.Answer{}).Select("count(*) as count, concat(@createTime, '-01') as create_time").
		Where("?<=created_at", monthlyBegin).Group("@createTime := DATE_FORMAT(created_at, '%Y-%m')").Scan(&counts).Error; err != nil {
		return nil, err
	}
	if err := calculateChartCounts(monthlyBegin, 1, 0, counts, resp.Chart.Monthly.NewAnswer, nil); err != nil {
		return nil, err
	}

	// =================================================================================================================
	//                                               total progress
	// =================================================================================================================
	var paper model.Paper
	if err := db.GetDB().Model(&model.Paper{}).Select("id, words").Where("inuse=?", true).Last(&paper).Error; err != nil {
		return nil, err
	}
	if err := db.GetDB().Model(&model.Answer{}).Where("paper_id=?", paper.Id).Count(&count).Error; err != nil {
		return nil, err
	}
	resp.TotalProgress = float32(count) / float32(len(regexp.MustCompile("[ \t\n]").Split(paper.Words, -1))*resp.TotalUsers)

	return resp, nil
}

func calculateChartCounts(begin time.Time, offsetMonth, offsetDay int, counts []*countTool, result []int, labels []string) error {
	for i, j := 0, 0; i < len(result); i++ {
		date := begin.AddDate(0, offsetMonth*i, offsetDay*i)
		if len(labels) > 0 {
			if offsetDay > 0 {
				labels[i] = date.Format("2 Jan")
			} else {
				labels[i] = date.Format("Jan")
			}
		}
		if j > len(counts)-1 {
			continue
		}
		if counts[j].CreateTime == "" {
			continue
		}
		parsed, err := time.ParseInLocation("2006-01-02", counts[j].CreateTime, time.Local)
		if err != nil {
			return err
		}
		for date.After(parsed) {
			j++
			if j > len(counts)-1 {
				break
			}
			parsed, err = time.ParseInLocation("2006-01-02", counts[j].CreateTime, time.Local)
			if err != nil {
				return err
			}
		}
		if date.Equal(parsed) {
			result[i] = counts[j].Count
		}
	}
	return nil
}
