package logger

import (
	"context"
	"fmt"
	"log"
	"os"
)

type Severity int

const (
	DEBUG Severity = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

func handleLog(ctx context.Context, severity Severity, message string) {
	prefix := "UNKNOWN"
	switch severity {
	case DEBUG:
		prefix = "DEBUG"
	case INFO:
		prefix = "INFO"
	case WARNING:
		prefix = "WARN"
	case ERROR:
		prefix = "ERROR"
	case CRITICAL:
		prefix = "FATAL"
	}
	log.Printf("%s: %s", prefix, message)
}

func Debugf(ctx context.Context, format string, v ...any) {
	handleLog(ctx, DEBUG, fmt.Sprintf(format, v...))
}

func Infof(ctx context.Context, format string, v ...any) {
	handleLog(ctx, INFO, fmt.Sprintf(format, v...))
}

func Warnf(ctx context.Context, format string, v ...any) {
	handleLog(ctx, WARNING, fmt.Sprintf(format, v...))
}

func Errorf(ctx context.Context, format string, v ...any) {
	handleLog(ctx, ERROR, fmt.Sprintf(format, v...))
}

func Fatalf(ctx context.Context, format string, v ...any) {
	handleLog(ctx, CRITICAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}
