package log

import (
	"errors"
	"io"
	"log"
)

const (
	INFO int = iota
	DEBUG
	TRACE
	ERROR
)

var (
	ErrLevelBig = errors.New("the level set is less than the level of the log exporter")
)

type Logger interface {
	Level() int
	SetLevel(level int)
	Debug(message string) error
	Info(message string) error
	Trace(message string) error
	Error(message string) error
}

func NewLogger(io io.Writer,level int) Logger {
	lg := log.New(io,"",log.Ltime | log.Lshortfile)
	return &stdLogger{
		lg: lg,
		level: level,
	}
}

type stdLogger struct {
	lg *log.Logger
	level int
}

func (s *stdLogger) Level() int {
	return s.level
}

func (s *stdLogger) SetLevel(level int) {
	s.level = level
}

// Verify the current level. If it is lower than or higher than the set level, it will not be output
func (s *stdLogger) checkLevel(level int) bool {
	return s.level < level
}

func (s *stdLogger) Debug(message string) error {
	if !s.checkLevel(DEBUG) {
		return ErrLevelBig
	}
	return nil
}

func (s *stdLogger) Info(message string) error {
	if !s.checkLevel(INFO) {
		return ErrLevelBig
	}
	return nil
}

func (s *stdLogger) Trace(message string) error {
	if !s.checkLevel(TRACE) {
		return ErrLevelBig
	}
	return nil
}

func (s *stdLogger) Error(message string) error {
	if !s.checkLevel(ERROR) {
		return ErrLevelBig
	}
	return nil
}
