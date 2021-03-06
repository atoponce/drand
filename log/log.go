package log

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
	lvl "github.com/go-kit/kit/log/level"
)

// Logger is a interface that can log to different levels.
type Logger interface {
	Info(keyvals ...interface{})
	Debug(keyvals ...interface{})
	Warn(keyvals ...interface{})
	Error(keyvals ...interface{})
	Fatal(keyvals ...interface{})
	// With returns a new Logger that inserts the given key value pairs for each
	// statements at each levels
	With(keyvals ...interface{}) Logger
}

const (
	// LogNone forbids any log entries
	LogNone int = iota
	// LogInfo sets the logging verbosity to info
	LogInfo
	// LogDebug sets the logging verbosity to debug
	LogDebug
)

// DefaultLevel is the default level where statements are logged. Change the
// value of this variable before init() to change the level of the default
// logger.
var DefaultLevel = LogInfo

// DefaultLogger is the default logger that only statemetns at the
// default level. The level is set by DefaultLevel inside init()..
var DefaultLogger Logger

func init() {
	DefaultLogger = NewLogger(DefaultLevel)
}

type kitLogger struct {
	log.Logger
}

// NewLogger returns a kit logger that prints statements at the given level.
func NewLogger(level int) Logger {
	var opt lvl.Option
	switch level {
	case LogNone:
		opt = lvl.AllowNone()
	case LogInfo:
		opt = lvl.AllowInfo()
	case LogDebug:
		opt = lvl.AllowDebug()
	default:
		panic("unknown log level")
	}
	return NewKitLogger(opt)
}

// NewKitLoggerFrom returns a Logger out of a go-kit/kit/log logger interface. The
// caller can set the options that it needs to the logger first. By default, it
// wraps the logger with a SyncLogger so concurrent calls are synchronized.
func NewKitLoggerFrom(l log.Logger) Logger {
	return &kitLogger{log.NewSyncLogger(l)}
}

// NewKitLogger returns a Logger based on go-kit/kit/log default logger
// structure that outputs to stdout. You can pass in options to only allow
// certain levels. By default, it also includes the caller stack.
func NewKitLogger(opts ...lvl.Option) Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	for _, opt := range opts {
		logger = lvl.NewFilter(logger, opt)
	}
	timestamp := log.TimestampFormat(time.Now, time.RFC1123)
	logger = log.With(logger, "ts", timestamp)
	logger = log.With(logger, "call", log.Caller(8))
	return NewKitLoggerFrom(logger)
}

func (k *kitLogger) Info(kv ...interface{}) {
	lvl.Info(k.Logger).Log(kv...)
}

func (k *kitLogger) Debug(kv ...interface{}) {
	lvl.Debug(k.Logger).Log(kv...)
}

func (k *kitLogger) Warn(kv ...interface{}) {
	lvl.Warn(k.Logger).Log(kv...)
}

func (k *kitLogger) Error(kv ...interface{}) {
	lvl.Error(k.Logger).Log(kv...)
}

func (k *kitLogger) Fatal(kv ...interface{}) {
	lvl.Error(k.Logger).Log(kv...)
	os.Exit(1)
}

func (k *kitLogger) With(kv ...interface{}) Logger {
	newLogger := log.With(k.Logger, kv...)
	return NewKitLoggerFrom(newLogger)
}
