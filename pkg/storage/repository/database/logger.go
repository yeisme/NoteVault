package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// logxAdapter 实现 gorm.Logger 接口，将日志转发到 logx
type logxAdapter struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

// NewLogxAdapter 创建一个新的 logx 适配器
func NewLogxAdapter() logger.Interface {
	return &logxAdapter{
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}
}

// LogMode 设置日志级别
func (l *logxAdapter) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 打印信息日志
func (l *logxAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		logx.Infof(msg, data...)
	}
}

// Warn 打印警告日志
func (l *logxAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logx.Slowf(msg, data...)
	}
}

// Error 打印错误日志
func (l *logxAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logx.Errorf(msg, data...)
	}
}

// Trace 记录 SQL 执行情况
func (l *logxAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 构建日志消息
	logMsg := fmt.Sprintf("[%.3fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)

	// 根据错误和执行时间判断日志级别
	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		// SQL 执行错误
		logx.Errorf("%s [ERROR: %v]", logMsg, err)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		// 慢查询
		logx.Slowf("%s [SLOW]", logMsg)
	default:
		// 普通日志
		if l.LogLevel >= logger.Info {
			logx.Infof("%s", logMsg)
		}
	}
}
