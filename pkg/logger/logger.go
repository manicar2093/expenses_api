package logger

import (
	"log"
	"os"
)

func New() *log.Logger {
	return log.New(os.Stdout, "INFO: ", log.Lshortfile|log.Ldate|log.Ltime)
}
