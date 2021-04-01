package utils

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

func InitLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
		PadLevelText:  true,
	})
}

func DebagLog(ctx context.Context, fields Fields) {
	logrus.WithFields(logrus.Fields(fields)).
		Debug("[request_id: ", ctx.Value("request_id"), "]")
}

func WarnLog(ctx context.Context, fields Fields) {
	logrus.WithFields(logrus.Fields(fields)).
		Warn("[request_id: ", ctx.Value("request_id"), "]")
}

func InfoLog(ctx context.Context, fields Fields) {
	logrus.WithFields(logrus.Fields(fields)).
		Info("[request_id: ", ctx.Value("request_id"), "]")
}
