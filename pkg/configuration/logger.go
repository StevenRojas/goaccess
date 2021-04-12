package configuration

import (
	"log/syslog"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// LoggerWrapper interface of logger wrapper
type LoggerWrapper interface {
	// Debug logger
	Debug(message string, metadata ...interface{})
	// Info logger
	Info(message string, metadata ...interface{})
	// Error logger
	Error(message string, metadata ...interface{})
}

type loggerWrapper struct {
	logger log.Logger
}

// NewLogger creates a new logger instance
func NewLogger(config ServerConfig) LoggerWrapper {
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(
		logger,
		//"at", log.DefaultTimestampUTC,
		//"caller", log.Caller(6),
	)
	switch config.LogLevel {
	case syslog.LOG_EMERG, syslog.LOG_CRIT, syslog.LOG_ALERT, syslog.LOG_ERR:
		logger = level.NewFilter(logger, level.AllowError())
	case syslog.LOG_WARNING:
		logger = level.NewFilter(logger, level.AllowWarn())
	case syslog.LOG_INFO, syslog.LOG_NOTICE:
		logger = level.NewFilter(logger, level.AllowInfo())
	case syslog.LOG_DEBUG:
		logger = level.NewFilter(logger, level.AllowDebug())
	}
	return &loggerWrapper{
		logger: logger,
	}
}

// Debug logger
func (l loggerWrapper) Debug(message string, metadata ...interface{}) {
	if len(metadata) > 0 {
		level.Debug(l.logger).Log("message", message, "meta", metadata)
	} else {
		level.Debug(l.logger).Log("message", message)
	}
}

// Info logger
func (l loggerWrapper) Info(message string, metadata ...interface{}) {
	if len(metadata) > 0 {
		level.Info(l.logger).Log("message", message, "meta", metadata)
	} else {
		level.Info(l.logger).Log("message", message)
	}
}

// Error logger
func (l loggerWrapper) Error(message string, metadata ...interface{}) {
	if len(metadata) > 0 {
		level.Error(l.logger).Log("message", message, "meta", metadata)
	} else {
		level.Error(l.logger).Log("message", message)
	}
}
