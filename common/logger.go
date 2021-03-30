package common

import "go.uber.org/zap"

var Log *zap.Logger

func InitLogger() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Log.Info("zap logger initiated")
}
