package librarian

import (
	"log"
	"os"
)

// NewLogger creates a new logger instance writing to the specified file path.
func NewLogger(filePath string) (*log.Logger, error) {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return log.New(logFile, "LIBRARIAN: ", log.Ldate|log.Ltime|log.Lshortfile), nil
}
