package model

// 用户注册信息
type User struct {
	Base
	Email            string `gorm:"size:100" json:"email"`
	Name             string `gorm:"size:255" json:"name"`
	ChineseClass     string `gorm:"size:50" json:"chinese_class"`
	HksLevel         string `gorm:"size:30" json:"hks_level"`
	EthnicBackground string `gorm:"size:50" json:"ethnic_background"`
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
