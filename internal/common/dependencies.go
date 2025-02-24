package common

import (
	"log"
	"os"
)

type Logger interface {
	Println(string)
}

type logger struct {
	logger *log.Logger
}

func NewLogger() Logger {
	lg := &logger{}
	lg.logger = log.New(os.Stdout, "twitter:", log.LstdFlags)
	return lg
}

func (l *logger) Println(str string) {
	l.logger.Println(str)
}
