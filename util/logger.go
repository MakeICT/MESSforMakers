package util

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const FILE_ONLY = "fileonly"
const FILE_AND_STDOUT = "fileandstdout"
const STDOUT_ONLY = "sdtoutonly"

const DEBUG = "debug"
const INFO = "info"
const WARN = "warn"
const ERROR = "error"
const FATAL = "fatal"

type Logger struct {
	*log.Logger
	file  *os.File
	level string
}

func NewLogger(logfile string, level string) (*Logger, error) {

	//Setup file logging
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %#v", err)
	}
	multiWritter := io.MultiWriter(f, os.Stdout)

	logger := log.New(multiWritter, "", log.Ldate|log.Ltime|log.Lshortfile)

	myLogger := &Logger{logger, f, DEBUG}

	err = myLogger.SetLevel(level)
	return myLogger, err
}

func (l *Logger) Close() {
	l.file.Close()
}

func (l *Logger) SetLevel(lev string) error {
	if lev == DEBUG || lev == INFO || lev == WARN || lev == ERROR || lev == FATAL {
		l.level = lev
		return nil
	}
	return errors.New("Level not recognized")
}

func (l *Logger) Warn(s string) {
	if l.level == WARN || l.level == INFO || l.level == DEBUG {
		l.Print(s)
	}
}
