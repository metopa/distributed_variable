package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/metopa/distributed_variable/common"
)

type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	fatal *log.Logger
}

var defaultLogger = Logger{
	info:  log.New(os.Stdout, "INFO  ", log.Ltime|log.Lmicroseconds),
	warn:  log.New(os.Stdout, "WARN  ", log.Ltime|log.Lmicroseconds),
	fatal: log.New(os.Stdout, "FATAL ", log.Ltime|log.Lmicroseconds|log.Lshortfile)}

func (l *Logger) Info(format string, v ...interface{}) {
	l.output(l.info, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.output(l.warn, format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.output(l.fatal, format, v...)
	os.Exit(1)
}

func (l *Logger) output(stream *log.Logger, format string, v ...interface{}) {
	stream.Output(4, fmt.Sprintf("%v "+format,
		append([]interface{}{common.PeekLogicalTimestamp()}, v...)...))
}

func Info(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Fatal(format string, v ...interface{}) {
	defaultLogger.Fatal(format, v...)
}
