package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

const (
	middlewareLevel       = "Middleware Level"
	usecaseLevel          = "Usecase Level"
	deliveryLevel         = "Delivery Level"
	repositoryLevel       = "Repository Level"
	responseLevel         = "Response Level"
	utilsLevel            = "Utils Level"
	serviceInterfaceLevel = "Service Interface Level"
)

type Fields map[string]interface{}

func InitLogger() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
		PadLevelText:  true,
	})
}

type EntryLog struct {
	level string
}

func (entry *EntryLog) DebugLog(ctx context.Context, fields Fields) {
	logrus.WithFields(logrus.Fields(fields)).
		Debug("[id: ", ctx.Value("request_id"), "] ", entry.level)
}

func (entry *EntryLog) WarnLog(ctx context.Context, fields Fields) {
	logrus.WithFields(logrus.Fields(fields)).
		Warn("[id: ", ctx.Value("request_id"), "] ", entry.level)
}

func (entry *EntryLog) InfoLog(ctx context.Context, fields Fields) {
	logrus.WithFields(logrus.Fields(fields)).
		Info("[id: ", ctx.Value("request_id"), "] ", entry.level)
}

func (entry *EntryLog) ErrorLog(ctx context.Context, err error) {
	logrus.WithFields(logrus.Fields{
		"error": err.Error(),
	}).Warn("[id: ", ctx.Value("request_id"), "] ", entry.level)
}

func (entry *EntryLog) InlineInfoLog(ctx context.Context, data interface{}) {
	logrus.WithFields(logrus.Fields{
		"info": data,
	}).Info("[id: ", ctx.Value("request_id"), "] ", entry.level)
}

func (entry *EntryLog) InlineDebugLog(ctx context.Context, data interface{}) {
	logrus.WithFields(logrus.Fields{
		"data": data,
	}).Debug("[id: ", ctx.Value("request_id"), "] ", entry.level)
}

func (entry *EntryLog) DataLog(ctx context.Context, name string, data interface{}) {
	logrus.WithFields(logrus.Fields{
		name: data,
	}).Info("[id: ", ctx.Value("request_id"), "] ", entry.level)
}

func MiddleLevel() *EntryLog {
	return &EntryLog{
		level: middlewareLevel,
	}
}

func DeliveryLevel() *EntryLog {
	return &EntryLog{
		level: deliveryLevel,
	}
}

func UsecaseLevel() *EntryLog {
	return &EntryLog{
		level: usecaseLevel,
	}
}

func RepoLevel() *EntryLog {
	return &EntryLog{
		level: repositoryLevel,
	}
}

func ResponseLevel() *EntryLog {
	return &EntryLog{
		level: responseLevel,
	}
}

func UtilsLevel() *EntryLog {
	return &EntryLog{
		level: utilsLevel,
	}
}

func ServiceInterfaceLevel() *EntryLog {
	return &EntryLog{
		level: serviceInterfaceLevel,
	}
}
