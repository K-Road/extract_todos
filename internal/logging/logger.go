package logging

import (
	"log"

	"github.com/natefinch/lumberjack"
)

var (
	WebServerLogger    *log.Logger
	WebServerLogWriter *lumberjack.Logger
	ExtractLogger      *log.Logger
	TUILogger          *log.Logger
)

func Init() {
	WebServerLogWriter = &lumberjack.Logger{
		Filename:   "./logs/webserver.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,  // number of backups
		MaxAge:     30, // days
		Compress:   true,
	}
	WebServerLogger = log.New(WebServerLogWriter, "WEBSERVER", log.LstdFlags|log.Lshortfile)

	ExtractLogger = log.New(&lumberjack.Logger{
		Filename:   "./logs/extract.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}, "EXTRACT", log.LstdFlags|log.Lshortfile)

	TUILogger = log.New(&lumberjack.Logger{
		Filename:   "./logs/tui.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}, "TUI", log.LstdFlags|log.Lshortfile)
}

func TUI() *log.Logger {
	return TUILogger
}
func WEB() *log.Logger {
	return WebServerLogger
}
func Extract() *log.Logger {
	return ExtractLogger
}
