package common

import (
	"go.uber.org/zap"
	"os"
	"read-test-server/common/gorm2/mysql"
)

const (
	AudioUploadRoot = "/asset/audio_upload/"
)

func InitAudioUploadRoot() {
	stat, err := os.Stat(AudioUploadRoot)
	if err != nil {
		if innerErr := os.MkdirAll(AudioUploadRoot, 0777); innerErr != nil {
			Log.Panic("InitAudioUploadRoot.MkdirAll", zap.Error(innerErr))
		}
		return
	}
	if !stat.IsDir() {
		Log.Panic("AudioUploadRoot is not Dir", zap.String("root", AudioUploadRoot))
	}
}

type ServerConfig struct {
	Mysql *mysql.DbConfig `json:"mysql"`
}