package middleware

import (
	"time"

	"github.com/borscht/backend/config"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	context *logrus.Entry
}

func (logger *LoggerMiddleware) Log(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		start := time.Now()
		result := next(c)

		logger.context.WithFields(logrus.Fields{
			"url":           c.Request().URL,
			"method":        c.Request().Method,
			"remote_addr":   c.Request().RemoteAddr,
			"work_time":     time.Since(start),
			"server_status": c.Response().Status,
		}).Info("request_id: ", c.Get("request_id"))

		return result
	}
}

func InitLoggerMiddleware() *LoggerMiddleware {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
		PadLevelText:  true,
	})
	logrus.WithFields(logrus.Fields{
		"host": config.Host,
		"port": config.ServerPort,
	}).Info("Starting server")

	contextLogger := logrus.WithFields(logrus.Fields{})

	return &LoggerMiddleware{
		context: contextLogger,
	}
}
