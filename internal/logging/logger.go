package logging

import (
	"log"
	"os"

	"github.com/natefinch/lumberjack"
)

var (
	WebServerLogger    *log.Logger
	WebServerLogWriter *lumberjack.Logger
	ExtractLogger      *log.Logger
	TUILogger          *log.Logger
	DBLogger           *log.Logger
	DATALogger         *log.Logger
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

	DBLogger = log.New(&lumberjack.Logger{
		Filename:   "./logs/db.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}, "DB", log.LstdFlags|log.Lshortfile)

	DATALogger = log.New(&lumberjack.Logger{
		Filename:   "./logs/data.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}, "DATA", log.LstdFlags|log.Lshortfile)
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
func DB() *log.Logger {
	return DBLogger
}
func DATA() *log.Logger {
	return DATALogger
}

func ExitWithError(log *log.Logger, msg string, err error) {
	log.Println(msg+":", err)
	os.Exit(1)
}
