package gormslog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Slog interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type Logger struct {
	logger                Slog
	SlowThreshold         time.Duration
	SkipErrRecordNotFound bool
	SkipErrContexCanceled bool
	Debug                 bool
	MsgFormatter          func(sql string, elapsed time.Duration, source string) string
}

func New(logger Slog) *Logger {
	if logger == nil {
		logger = slog.Default()
	}

	return &Logger{
		logger:                logger,
		SkipErrRecordNotFound: true,
		SkipErrContexCanceled: true,
		Debug:                 true,
		MsgFormatter: func(sql string, elapsed time.Duration, source string) string {
			return fmt.Sprintf("%s [%s] %s", sql, elapsed, source)
		},
	}
}

func (l *Logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.logger.InfoContext(ctx, fmt.Sprintf(s, args...))
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.logger.WarnContext(ctx, fmt.Sprintf(s, args...))
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.logger.ErrorContext(ctx, fmt.Sprintf(s, args...))
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	source := sourceShort(utils.FileWithLineNum())
	sql, _ := fc()
	msg := l.MsgFormatter(sql, elapsed, source)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound {
			return
		}

		if errors.Is(err, context.Canceled) && !l.SkipErrContexCanceled {
			return
		}

		l.logger.ErrorContext(ctx, msg, "error", err)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger.WarnContext(ctx, msg)
		return
	}

	if l.Debug {
		l.logger.DebugContext(ctx, msg)
	}
}

func sourceShort(s string) string {
	parts := strings.Split(s, "/")
	if len(parts) <= 2 {
		return s
	}

	return strings.Join(parts[len(parts)-2:], "/")
}
