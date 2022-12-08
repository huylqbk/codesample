package logger

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Fatal(msg string, keyvals ...interface{})
	SetCaller() Logger
	SetLevel(level int) Logger
	LogFile(path string) Logger
}

var logObj Logger

func Get() Logger {
	return logObj
}

func caller(depth int) string {
	_, file, line, _ := runtime.Caller(depth)
	idx := strings.LastIndexByte(file, '/')
	return fmt.Sprintf("%s:%d", file[idx+1:], line)
}

func getWriter(path string) io.Writer {
	if path == "" {
		path = "./logger"
	}
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s/%s.log", path, "%Y-%m-%d"),
		rotatelogs.WithMaxAge(time.Hour*24*10),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic("Failed to Initialize Log File")
	}
	return writer
}
