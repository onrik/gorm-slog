package gormslog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	logger                *slog.Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

func New(logger *slog.Logger) *Logger {
	if logger == nil {
		logger = slog.Default()
	}
	return &Logger{
		logger:                logger,
		SkipErrRecordNotFound: true,
		Debug:                 true,
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
	sql, _ := fc()
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		l.logger.With("error", err).ErrorContext(ctx, fmt.Sprintf("%s [%s]", sql, elapsed))
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger.WarnContext(ctx, fmt.Sprintf("%s [%s]", sql, elapsed))
		return
	}

	if l.Debug {
		l.logger.DebugContext(ctx, fmt.Sprintf("%s [%s]", sql, elapsed))
	}
}
