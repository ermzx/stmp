package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger 初始化日志工具
func InitLogger() {
	Logger = logrus.New()

	// 设置日志格式
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志输出
	Logger.SetOutput(os.Stdout)

	// 设置日志级别
	Logger.SetLevel(logrus.DebugLevel)
}

// Info 记录信息日志
func Info(args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Info(args...)
}

// Error 记录错误日志
func Error(args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Error(args...)
}

// Debug 记录调试日志
func Debug(args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Debug(args...)
}

// Warn 记录警告日志
func Warn(args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Warn(args...)
}

// Infof 记录格式化信息日志
func Infof(format string, args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Infof(format, args...)
}

// Errorf 记录格式化错误日志
func Errorf(format string, args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Errorf(format, args...)
}

// Debugf 记录格式化调试日志
func Debugf(format string, args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Debugf(format, args...)
}

// Warnf 记录格式化警告日志
func Warnf(format string, args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Warnf(format, args...)
}