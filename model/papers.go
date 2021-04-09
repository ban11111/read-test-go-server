package model

type Paper struct {
	Base
	Name     string `gorm:"size:50;COMMENT:'试卷名'" json:"name"`
	Version  int16  `gorm:"COMMENT:'试卷版本'" json:"version"`
	Words    string `gorm:"type:text;COMMENT:'字词题目, 格式: 字词[空格]字词[空格]字词'" json:"words"`
	Interval int    `gorm:"COMMENT:'做题间隔, 单位: 秒'" json:"interval"`
	Inuse    bool   `gorm:"COMMENT:'生效标识'" json:"inuse"`
}

type PaperSnapshot struct {
	Paper
	PaperId uint `json:"paper_id"`
}
