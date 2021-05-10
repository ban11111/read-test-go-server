package model

// 用户注册信息
type User struct {
	Base
	Email                  string `gorm:"size:100" json:"email" csv:"email" width:"30"`
	Name                   string `gorm:"size:255" json:"name" csv:"name" width:"18"`
	ChineseClass           string `gorm:"size:50" json:"chinese_class" csv:"chinese_class" width:"13"`
	HksLevel               string `gorm:"size:30" json:"hks_level" csv:"hks_level" width:"12.5"`
	EthnicBackground       string `gorm:"size:50" json:"ethnic_background" csv:"ethnic_background" width:"18"`
	HasChineseAcquaintance bool   `json:"has_chinese_acquaintance" csv:"cn_relation" width:"11"`
	AcquaintanceDetail     string `gorm:"size:255" json:"acquaintance_detail" csv:"acquaintance_detail" width:"25"`
	Score                  string `sql:"-" json:"-" csv:"score"`
}

// 返给前端, 包含试卷信息
type UserInfo struct {
	User
	Papers []*PaperUser `json:"papers" gorm:"foreignKey:Uid;references:Id"` // 做过的试卷
}

func (i UserInfo) TableName() string {
	return "users"
}

type PaperUser struct {
	Id    uint   `json:"id" gorm:"primaryKey"`
	Uid   uint   `json:"uid"`
	Pid   uint   `json:"pid"`
	PName string `json:"p_name" gorm:"size:50"`
}
