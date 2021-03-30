package common

import (
	dbConfig "github.com/dipperin/go-ms-toolkit/db-config"
	"github.com/dipperin/go-ms-toolkit/orm/gorm/mysql"
	"go.uber.org/zap"
	"os"
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
	Mysql *dbConfig.DbConfig `json:"mysql"`
}

func (s *ServerConfig) MakeDB() mysql.DB {
	return mysql.MakeDB(s.Mysql)
}