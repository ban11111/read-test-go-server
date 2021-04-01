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
