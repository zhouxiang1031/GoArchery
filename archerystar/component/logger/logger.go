package logger

import (
	"fmt"
	"runtime"
	"strings"

	"archerystar/config"

	"github.com/sirupsen/logrus"
)

type Fields logrus.Fields

// Log is the default logger
var logManager = logrus.New() //initLogger(logrus.DebugLevel)

func init() {
	InitLogger()
}

func InitLogger() {
	//logManager = logrus.New()
	logManager.Formatter = new(logrus.TextFormatter)
	logManager.Level = logrus.DebugLevel
	if config.Gameconfig().Server.LogFile {
		logManager.SetOutput(newWirter())
	}

	fmt.Println("init log ok!")

	Info("logmode", "init log ok!")
	/*
		log := plog.WithFields(logrus.Fields{
			"app": "archery",
		})
	*/
}

func SetLogLevel(level logrus.Level) {
	logManager.Level = level
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<not found file>"
		line = -1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func SetLogFormatter(formatter logrus.Formatter) {
	logManager.Formatter = formatter
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	//fmt.Println(file)
	//fmt.Println(line)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}

// Debug
func Debug(title string, format string, args ...interface{}) {
	if logManager.Level >= logrus.DebugLevel {
		entry := logManager.WithFields(logrus.Fields{"title": title})
		//entry.Data["file"] = fileInfo(2)
		entry.Debugf(format, args...)
	}
}

// 带有field的Debug
func DebugWithFields(l interface{}, f Fields) {
	if logManager.Level >= logrus.DebugLevel {
		entry := logManager.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Debug(l)
	}
}

// Info
func Info(title string, format string, args ...interface{}) {
	if logManager.Level >= logrus.InfoLevel {
		entry := logManager.WithFields(logrus.Fields{"title": title})
		//entry.Data["file"] = fileInfo(2)
		entry.Infof(format, args...)
	}
}

// 带有field的Info
func InfoWithFields(l interface{}, f Fields) {
	if logManager.Level >= logrus.InfoLevel {
		entry := logManager.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Info(l)
	}
}

// Warn
func Warn(title string, format string, args ...interface{}) {
	if logManager.Level >= logrus.WarnLevel {
		entry := logManager.WithFields(logrus.Fields{"title": title})
		entry.Data["file"] = fileInfo(2)
		entry.Warnf(format, args...)
	}
}

// 带有Field的Warn
func WarnWithFields(l interface{}, f Fields) {
	if logManager.Level >= logrus.WarnLevel {
		entry := logManager.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Warn(l)
	}
}

// Error
func Error(title string, format string, args ...interface{}) {
	if logManager.Level >= logrus.ErrorLevel {
		entry := logManager.WithFields(logrus.Fields{"title": title})
		entry.Data["file"] = fileInfo(2)
		//entry.Error(args...)
		entry.Errorf(format, args...)
	}
}

// 带有Fields的Error
func ErrorWithFields(l interface{}, f Fields) {
	if logManager.Level >= logrus.ErrorLevel {
		entry := logManager.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Error(l)
	}
}

// Fatal
func Fatal(title string, format string, args ...interface{}) {
	if logManager.Level >= logrus.FatalLevel {
		entry := logManager.WithFields(logrus.Fields{"title": title})
		entry.Data["file"] = fileInfo(2)
		entry.Fatalf(format, args...)
	}
}

// 带有Field的Fatal
func FatalWithFields(l interface{}, f Fields) {
	if logManager.Level >= logrus.FatalLevel {
		entry := logManager.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(l)
	}
}

// Panic
func Panic(title string, format string, args ...interface{}) {
	if logManager.Level >= logrus.PanicLevel {
		entry := logManager.WithFields(logrus.Fields{"title": title})
		entry.Data["file"] = fileInfo(2)
		entry.Panicf(format, args...)
	}
}

// 带有Field的Panic
func PanicWithFields(l interface{}, f Fields) {
	if logManager.Level >= logrus.PanicLevel {
		entry := logManager.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Panic(l)
	}
}
