package database

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Logger struct {
	log                   *logrus.Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

func NewLogger(l *logrus.Logger) *Logger {
	return &Logger{
		log:                   l,
		SkipErrRecordNotFound: true,
	}
}

func (l *Logger) LogMode(gormlog.LogLevel) gormlog.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.log.WithField("component", "database").WithContext(ctx).Infof(s, args...)
}

func (l *Logger) Debug(ctx context.Context, s string, args ...interface{}) {
	l.log.WithField("component", "database").WithContext(ctx).Tracef(s, args...) // Re-route to trace as we don't want SQL queries in debug output
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.log.WithField("component", "database").WithContext(ctx).Warnf(s, args...)
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.log.WithField("component", "database").WithContext(ctx).Errorf(s, args...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{
		"component": "database",
	}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.log.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.log.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	l.log.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
}
