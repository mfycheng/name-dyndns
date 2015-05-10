// Package log provides a basic global logger for name-dyndns.
package log

import (
	"io"
	"log"
)

// Global logger.
var Logger *log.Logger

// Init intializes the logger with a specific io.Writer.
// This function is generally called near startup.
func Init(writer io.Writer) {
	Logger = log.New(writer, "", log.Ldate|log.Ltime)
}
