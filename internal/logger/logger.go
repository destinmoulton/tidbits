package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tidbits/internal/cliflags"
	"tidbits/internal/utils"
)

type Logger struct {
	logfileHandler *os.File
	logger         *log.Logger
}

func NewLogger(flags *cliflags.FlagOptions) *Logger {
	var h *os.File = os.Stdout
	if flags.LogToFile {
		// create a file handler
		dir := utils.GetConfigDir()

		path := filepath.Join(dir, flags.LogFile)
		var err error

		h, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			fmt.Printf("failed to open log file: %s\n", path)
			os.Exit(2)
		}
	}
	l := log.New(h, "", log.LstdFlags)

	return &Logger{
		logfileHandler: h,
		logger:         l,
	}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.SetPrefix("DEBUG:")
	l.logger.Println(msg, args)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.SetPrefix("WARN:")
	l.logger.Println(msg, args)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.SetPrefix("ERROR:")
	l.logger.Println(msg, args)
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.logger.Fatalln(msg, args)
}

func (l *Logger) Close() {
	l.logfileHandler.Close()
	l.Close()
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.SetPrefix("INFO:")
	l.logger.Println(msg, args)
}
