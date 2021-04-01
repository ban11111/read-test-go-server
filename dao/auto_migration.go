package dao

import (
	"go.uber.org/zap"
	"read-test-server/common"
	"read-test-server/model"
)

func AutoMigration() {
	if err := db.GetDB().AutoMigrate(&model.User{}); err != nil {
		common.Log.Panic("AutoMigration.User{}", zap.Error(err))
	}
	if err := db.GetDB().AutoMigrate(&model.Paper{}); err != nil {
		common.Log.Panic("AutoMigration.Paper{}", zap.Error(err))
	}
	if err := db.GetDB().AutoMigrate(&model.Answer{}); err != nil {
		common.Log.Panic("AutoMigration.Answer{}", zap.Error(err))
	}
}