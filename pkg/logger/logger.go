package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func New() *log.Logger {
	return log.New(os.Stdout, "INFO: ", log.Lshortfile|log.Ldate|log.Ltime)
}

func NewWithPrefix(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("%s INFO: ", strings.ToUpper(prefix)), log.Lshortfile|log.Ldate|log.Ltime)
}
