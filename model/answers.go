package model

import "time"

// gorm 通用前置字段
type Base struct {
	Id        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at" gorm:"COMMENT:'创建时间'"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"COMMENT:'更新时间'"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"COMMENT:'软删除时间'"`
}

// 定义用户答题结果的 数据库表
type Answer struct {
	Base
	PaperId     uint   `gorm:"index:idx_paper_result,priority:2" json:"paper_id"` // 关联 Paper表 id
	Uid         uint   `gorm:"index:idx_paper_result,priority:1" json:"uid"`
	WordIndex   int    `gorm:"index:idx_paper_result,priority:3" json:"word_index"`
	Word        string `gorm:"size:50" json:"word"`
	AudioUrl    string `gorm:"size:255" json:"audio_url"`
	Translation string `gorm:"size:255" json:"translation"`
	Duration    int    `json:"duration"` // 该题耗时, 单位: 毫秒
	Device      string `gorm:"size:50" json:"device"`
}
