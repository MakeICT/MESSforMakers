package util

import (
	"fmt"
	"io"
	"log"
	"os"
)

//Logger extends the standard logger to add multi-logging and log levels
type Logger struct {
	*log.Logger
	file *os.File
}

//NewLogger sets up a new logger based on settings provided
func NewLogger() (*Logger, error) {

	//Setup file logging
	//TODO add settings
	f, err := os.OpenFile("makeict.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v", err)
	}
	multiWritter := io.MultiWriter(f, os.Stdout)

	logger := log.New(multiWritter, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{logger, f}, nil
}

//Close should be deferred any time writing to a file is selected
func (l *Logger) Close() {
	l.file.Close()
}
