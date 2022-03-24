package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Logger struct {
	level  string
	writer *bufio.Writer
}

func New(level string, fileName string) *Logger {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Println(err)
	}

	return &Logger{level: level, writer: bufio.NewWriter(file)}
}

func (l Logger) Debug(msg string) {
	if l.level == "debug" {
		fmt.Println(msg)
		_, _ = l.writer.WriteString(msg + "\n")
		l.writer.Flush()
	}
}

func (l Logger) Info(msg string) {
	if l.level == "info" {
		fmt.Println(msg)
		_, _ = l.writer.WriteString(msg + "\n")
		l.writer.Flush()
	}
}

func (l Logger) Warning(msg string) {
	if l.level == "warning" {
		fmt.Println(msg)
		_, _ = l.writer.WriteString(msg + "\n")
		l.writer.Flush()
	}
}

func (l Logger) Error(msg string) {
	if l.level == "error" {
		_, _ = l.writer.WriteString(msg + "\n")
		l.writer.Flush()
	}
}
