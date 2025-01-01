package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

var logLevel LogLevel = LogLevelFatal

var moduleName string = ""

func Init(module string) {
	if module != "" {
		moduleName = module
	}

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		logLevel = LogLevelDebug
	case "info":
		logLevel = LogLevelInfo
	case "warning":
		logLevel = LogLevelWarn
	case "error":
		logLevel = LogLevelError
	case "fatal":
		logLevel = LogLevelFatal
	}

	log.SetFlags(0)
}

func formatLog(prefix string, v ...any) string {
	_, file, line, ok := runtime.Caller(2) // Skip two frames to get the original caller
	if !ok {
		file = "???"
		line = 0
	} else if moduleName != "" {
		if idx := strings.Index(file, moduleName); idx != -1 {
			file = file[(idx + len(moduleName) + 1):]
		}
	}

	// Format the timestamp, file, line, and the actual message
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	return fmt.Sprintf("%s %s %s:%d\n%s", prefix, timestamp, file, line, fmt.Sprintln(v...))
}

func Debug(v ...any) {
	if LogLevelDebug >= logLevel {
		log.Print(formatLog("\033[47mDEBUG\033[0m", v...))
	}
}

func Info(v ...any) {
	if LogLevelInfo >= logLevel {
		log.Print(formatLog("\033[106mINFO \033[0m", v...))
	}
}

func Warn(v ...any) {
	if LogLevelWarn >= logLevel {
		log.Print(formatLog("\033[103mWARN \033[0m", v...))
	}
}

func Error(v ...any) {
	if LogLevelError >= logLevel {
		log.Print(formatLog("\033[101mERROR\033[0m", v...))
	}
}

func Fatal(v ...any) {
	if LogLevelFatal >= logLevel {
		log.Fatal(formatLog("\033[105mFATAL\033[0m", v...))
	}
}
