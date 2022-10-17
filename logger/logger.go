package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	log    *logrus.Logger
	caller bool
	level  int
	file   bool
}

func NewLogger(sentryUrl string) Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	if sentryUrl != "" {
		hook, err := logrus_sentry.NewSentryHook(sentryUrl, []logrus.Level{
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

func (l *logrusLogger) LogFile() Logger {
	l.file = true
	mw := io.MultiWriter(os.Stdout, getWriter(""))
	l.log.SetOutput(mw)
	return l
}

func (l *logrusLogger) SetLevel(level int) Logger {
	l.level = level
	l.log.SetLevel(logrus.Level(level))
	return l
}

func (l *logrusLogger) SetCaller() Logger {
	l.caller = true
	return l
}

func (l *logrusLogger) Debug(msg string, keyvals ...interface{}) {
	l.log.WithFields(l.append(keyvals...)).Debug(msg)
}

func (l *logrusLogger) Info(msg string, keyvals ...interface{}) {
	l.log.WithFields(l.append(keyvals...)).Info(msg)
}

func (l *logrusLogger) Warn(msg string, keyvals ...interface{}) {
	l.log.WithFields(l.append(keyvals...)).Warn(msg)
}

func (l *logrusLogger) Error(msg string, keyvals ...interface{}) {
	l.log.WithFields(l.append(keyvals...)).Error(msg)
}

func (l *logrusLogger) Fatal(msg string, keyvals ...interface{}) {
	l.log.WithFields(l.append(keyvals...)).Fatal(msg)
}

func (l *logrusLogger) append(keyvals ...interface{}) logrus.Fields {
	fields := make(logrus.Fields)
	if l.caller {
		fields["caller"] = caller(3)
	}
	if len(keyvals) <= 1 {
		fields["data"] = keyvals
		return fields
	}
	len := len(keyvals)
	if len%2 == 0 {
		for i := 0; i < len; i += 2 {
			fields[fmt.Sprintf("%v", keyvals[i])] = keyvals[i+1]
		}
	} else {
		for i := 0; i < len-1; i += 2 {
			fields[fmt.Sprintf("%v", keyvals[i])] = keyvals[i+1]
		}
		fields["_extra"] = keyvals[len-1]
	}

	return fields
}
