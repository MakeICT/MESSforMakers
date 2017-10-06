package util

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	file *os.File
}

func NewLogger() (*Logger, error) {

	//Setup file logging
	f, err := os.OpenFile("makeict.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v", err)
	}
	multiWritter := io.MultiWriter(f, os.Stdout)

	logger := log.New(multiWritter, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{logger, f}, nil
}

func (l *Logger) Close() {
	l.file.Close()
}
