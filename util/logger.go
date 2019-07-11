package util

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	file        *os.File
	DumpRequest bool
}

func NewLogger(dr bool) (*Logger, error) {

	//Setup file logging
	//TODO make the output file a parameter that can be defined in config.json
	f, err := os.OpenFile("makeict.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v", err)
	}
	multiWritter := io.MultiWriter(f, os.Stdout)

	logger := log.New(multiWritter, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{logger, f, dr}, nil
}

func (l *Logger) Close() {
	l.file.Close()
}
