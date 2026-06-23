package logger

import "go.uber.org/zap"

var Instance *zap.SugaredLogger

func InitializeLogger() error {
	cfg, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	Instance = cfg.Sugar()
	return nil
}
