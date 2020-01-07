package util

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

//Constants for accepting parameters indicating where logged messages should be output. Currently unused.
const (
	FileOnly      = "fileonly"
	FileAndStdout = "fileandstdout"
	StdoutOnly    = "sdtoutonly"
)

//Flags used for selecting the info-only logging or debug level.
const (
	DEBUG = 0
	INFO  = 1
)

// Logger adds log leves and multi logging
type Logger struct {
	*log.Logger
	file        *os.File
	level       int
	DumpRequest bool
}

// NewLogger takes a string and some options and initializes a logger
func NewLogger(logfile string, dr bool, level int) (*Logger, error) {

	//Setup file logging
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v\n", err)
	}
	multiWritter := io.MultiWriter(f, os.Stdout)

	logger := Logger{
		log.New(multiWritter, "", log.Ldate|log.Ltime|log.Lshortfile),
		f,
		DEBUG,
		dr,
	}

	err = logger.SetLevel(level)

	return &logger, err
}

//Close should be deferred any time writing to a file is selected
func (l *Logger) Close() {
	l.file.Close()
}

//SetLevel allows changing the log level after a logger has been created.
func (l *Logger) SetLevel(lev int) error {
	if lev == DEBUG || lev == INFO {
		l.level = lev
		return nil
	}
	return errors.New("level not recognized")
}

//Debug extends fmt.Print to only log messages if the log level is set to DEBUG
func (l *Logger) Debug(s ...interface{}) {
	if l.level <= DEBUG {
		l.SetPrefix("DEBUG: ")
		l.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		l.Print(s...)
	}
}

//Debugf extends fmt.Printf to only log messages if the log level is set to DEBUG
func (l *Logger) Debugf(s string, data ...interface{}) {
	if l.level <= DEBUG {
		l.SetPrefix("DEBUG: ")
		l.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		l.Printf(s, data...)
	}
}

//Info extends fmt.Print
func (l *Logger) Info(s ...interface{}) {
	if l.level <= INFO {
		l.SetPrefix("INFO: ")
		l.SetFlags(log.Ldate | log.Ltime)
		l.Print(s...)
	}
}

//Infof extends fmt.Printf
func (l *Logger) Infof(format string, data ...interface{}) {
	if l.level <= INFO {
		l.SetPrefix("INFO: ")
		l.SetFlags(log.Ldate | log.Ltime)
		l.Printf(format, data...)
	}
}
