package logadapter

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// GormLogger model
type GormLogger struct {
	*Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

// NewGormLogger new gorm logger
func NewGormLogger() *GormLogger {
	l.SetFormatter(PrettyJSONFormat)
	return &GormLogger{
		SkipErrRecordNotFound: true,
		Debug:                 true,
		Logger:                l,
		SlowThreshold:         time.Second,
		SourceField:           DefaultSourceField,
	}
}

// LogMode get log mode
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	if level == gormlogger.Silent {
		l.Debug = false
	} else {
		l.Debug = true
	}

	return l
}

// Info log infor
func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Infof(s, args...)
}

// Warn log warn
func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Warnf(s, args...)
}

// Error log error
func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Errorf(s, args...)
}

// Trace log sql trace
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, row := fc()
	fields := logrus.Fields{
		"type":       LogTypeSQL,
		"row":        row,
		"latency_ms": elapsed.Milliseconds(),
		"latency":    elapsed.String(),
	}

	if l.Debug {
		fields["query"] = sqlMask(sql)
	}

	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.Logger.WithContext(ctx).WithFields(fields).Error()
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logger.WithContext(ctx).WithFields(fields).Warn()
		return
	}

	if l.Debug {
		fmt.Println("============ GORM log ===========")
		l.Logger.WithContext(ctx).WithFields(fields).Debug()
	}
}

func sqlMask(sql string) string {
	if len(sql) < DefaultLargeFieldLength {
		return sql
	}

	largeFields := []string{`<!DOCTYPE html PUBLIC`, `\"image\":`}

	// Mask large fields
	for i := 0; i < len(largeFields); i++ {
		index := strings.Index(sql, largeFields[i])
		if index != -1 {
			sql = sql[:index+50] + " ... " + sql[len(sql)-50:]
		}
	}

	// Mask json fields
	re := regexp.MustCompile(`\{.*\}`)
	match := re.FindStringIndex(sql)
	if len(match) > 2 && match[len(match)-1]-match[0] > DefaultLargeFieldLength {
		sql = sql[:match[0]+50] + " ... " + sql[match[len(match)-1]-50:]
	}

	return sql
}
