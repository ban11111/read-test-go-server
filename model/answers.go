package model

import "time"

// gorm 通用前置字段
type Base struct {
	Id        uint       `json:"id" gorm:"primaryKey" csv:"id" width:"8"`
	CreatedAt time.Time  `json:"created_at" gorm:"COMMENT:'创建时间';index" csv:"create_time" width:"18"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"COMMENT:'更新时间'" csv:"-"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"COMMENT:'软删除时间'" csv:"-"`
}

// 定义用户答题结果的 数据库表
type Answer struct {
	Base
	PaperId     uint   `gorm:"index:idx_paper_result,priority:2" json:"paper_id" csv:"paper_id" width:"8.86"` // 关联 Paper表 id
	Uid         uint   `gorm:"index:idx_paper_result,priority:1" json:"uid" csv:"uid" width:"6.14"`
	WordIndex   int    `gorm:"index:idx_paper_result,priority:3" json:"word_index" csv:"word_index" width:"11"`
	Word        string `gorm:"size:50" json:"word" csv:"word" width:"7"`
	AudioUrl    string `gorm:"size:255" json:"audio_url" csv:"audio_url" width:"30"`
	Translation string `gorm:"size:255" json:"translation" csv:"translation" width:"13.14"`
	Duration    int    `json:"duration" csv:"duration(ms)" width:"13"` // 该题耗时, 单位: 毫秒
	Device      string `gorm:"size:50" json:"device" csv:"device" width:"11.5"`
}
