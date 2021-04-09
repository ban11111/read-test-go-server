package model

const (
	SettingBitRate    = "BitRate"
	SettingSampleRate = "SampleRate"
)

type GlobalSetting struct {
	Key   string `gorm:"size:50" json:"key"`
	Value string `gorm:"size:100" json:"value"`
}