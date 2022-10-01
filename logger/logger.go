package logger

import (
	"fmt"
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(keyvals ...interface{})
	Info(keyvals ...interface{})
	Infof(msg string, keyvals ...interface{})
	Warn(keyvals ...interface{})
	Error(keyvals ...interface{})
	Errorf(msg string, keyvals ...interface{})
}

var logObj Logger

type logrusLogger struct {
	log *logrus.Logger
}

func Get() Logger {
	if logObj == nil {
		NewLogger("")
	}
	return logObj
}

func NewLogger(sentryUrl string) Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	if sentryUrl != "" {
		hook, err := logrus_sentry.NewSentryHook(sentryUrl, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})

		if err == nil {
			log.Hooks.Add(hook)
		}
	}

	logObj = &logrusLogger{log: log}
	return logObj
}

func (l *logrusLogger) Debug(keyvals ...interface{}) {
	l.log.Debug(keyvals...)
}

func (l *logrusLogger) Info(keyvals ...interface{}) {
	l.log.Info(keyvals...)
}

func (l *logrusLogger) Infof(msg string, keyvals ...interface{}) {
	field := make(logrus.Fields)
	for i := 0; i < len(keyvals); i += 2 {
		field[fmt.Sprintf("%v", keyvals[i])] = keyvals[i+1]
	}
	l.log.WithFields(field).Info(msg)
}

func (l *logrusLogger) Warn(keyvals ...interface{}) {
	l.log.Warn(keyvals...)
}

func (l *logrusLogger) Error(keyvals ...interface{}) {
	l.log.Error(keyvals...)
}

func (l *logrusLogger) Errorf(msg string, keyvals ...interface{}) {
	field := make(logrus.Fields)
	for i := 0; i < len(keyvals); i += 2 {
		field[fmt.Sprintf("%v", keyvals[i])] = keyvals[i+1]
	}
	l.log.WithFields(field).Error(msg)
}
