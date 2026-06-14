package util

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	fileLogger *log.Logger
	consoleLogger *log.Logger
}

func NewLogger(logPath string) *Logger {
	logFileName := filepath.Join(logPath, time.Now().Format("2006-01-02")+".log")
	
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		file = os.Stderr
	}

	return &Logger{
		fileLogger: log.New(file, "", log.LstdFlags|log.Lmicroseconds),
		consoleLogger: log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds),
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.fileLogger.Print(append([]interface{}{"[INFO]"}, v...)...)
	l.consoleLogger.Print(append([]interface{}{"[INFO]"}, v...)...)
}

func (l *Logger) Error(v ...interface{}) {
	l.fileLogger.Print(append([]interface{}{"[ERROR]"}, v...)...)
	l.consoleLogger.Print(append([]interface{}{"[ERROR]"}, v...)...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.fileLogger.Print(append([]interface{}{"[WARN]"}, v...)...)
	l.consoleLogger.Print(append([]interface{}{"[WARN]"}, v...)...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.fileLogger.Print(append([]interface{}{"[DEBUG]"}, v...)...)
	l.consoleLogger.Print(append([]interface{}{"[DEBUG]"}, v...)...)
}
