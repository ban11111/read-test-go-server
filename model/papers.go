package model

// todo, 定义试卷, 需要的参数:
type Paper struct {
	Base
	Name     string `gorm:"size:50;COMMENT:'试卷名'" json:"name"`
	Version  string `gorm:"size:10;COMMENT:'试卷版本'" json:"version"`
	Words    string `gorm:"type:text;COMMENT:'字词题目, 格式: 字词[空格]字词[空格]字词'" json:"words"`
	Interval int    `gorm:"COMMENT:'做题间隔, 单位: 秒'" json:"interval"`
}
