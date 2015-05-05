package log

import (
	"io"
	"log"
)

var Logger *log.Logger

func Init(writer io.Writer) {
	Logger = log.New(writer, "", log.Ldate|log.Ltime)
}
