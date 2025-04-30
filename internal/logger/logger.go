package logger

import (
	"context"
	"fmt"
	"log"
	"os"
)

type Severity int

const (
	Debug Severity = iota
	Info
	Warning
	Error
	Critical
)

func handleLog(ctx context.Context, severity Severity, message string) {
	prefix := "UNKNOWN"
	switch severity {
	case Debug:
		prefix = "DEBUG"
	case Info:
		prefix = "INFO"
	case Warning:
		prefix = "WARN"
	case Error:
		prefix = "ERROR"
	case Critical:
		prefix = "FATAL"
	}
	log.Printf("%s: %s", prefix, message)
}

func Debugf(ctx context.Context, format string, v ...any) {
	handleLog(ctx, Debug, fmt.Sprintf(format, v...))
}

func Infof(ctx context.Context, format string, v ...any) {
	handleLog(ctx, Info, fmt.Sprintf(format, v...))
}

func Warnf(ctx context.Context, format string, v ...any) {
	handleLog(ctx, Warning, fmt.Sprintf(format, v...))
}

func Errorf(ctx context.Context, format string, v ...any) {
	handleLog(ctx, Error, fmt.Sprintf(format, v...))
}

func Fatalf(ctx context.Context, format string, v ...any) {
	handleLog(ctx, Critical, fmt.Sprintf(format, v...))
	os.Exit(1)
}
