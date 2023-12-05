package logadapter

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// custom logtype
const (
	LogTypeAPI                = "api"
	LogTypeRequest            = "request"
	LogTypeResponse           = "response"
	LogTypeError              = "error"
	LogTypeDebug              = "debug"
	LogTypeWarn               = "warn"
	LogTypeInfo               = "info"
	LogTypeSQL                = "sql"
	LogTypeTrace              = "trace"
	LogTypeHTTPClientRequest  = "http-client-request"
	LogTypeHTTPClientResponse = "http-client-response"
)

// custom constants
const (
	DefaultTimestampFormat  = "2006-01-02 15:04:05.00000"
	DefaultSourceField      = "source"
	DefaultLargeFieldLength = 1000
)

// Key is key context type
type Key string

// Exported constanst
const (
	CorrelationIDKey Key = "X-User-Correlation-Id"
	RequestIDKey     Key = "X-Request-ID"
	UserInfoKey      Key = "X-User-Info"
	MessageIDKey     Key = "X-Message-ID"
)

// LogFormat log format
type LogFormat uint32

// custom log format
const (
	JSONFormat LogFormat = iota
	PrettyJSONFormat
	TextFormat
)

// Level log level
type Level uint32

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// Config config instance log
type Config struct {
	LogLevel        Level
	LogFormat       LogFormat
	TimestampFormat string // if empty, use default timestamp format
}

// Logger instance
type Logger struct {
	*log.Logger
	logFormat       LogFormat
	timestampFormat string
}

var baseSourceDir string

var l *Logger

func init() {
	l = New()
	sourceDir()
}

// SetFormatter set logger formatter
func SetFormatter(logFormat LogFormat) { l.SetFormatter(logFormat) }

// SetFormatter set logger formatter
func (l *Logger) SetFormatter(logFormat LogFormat) {
	switch logFormat {
	case JSONFormat:
		l.Logger.SetFormatter(&log.JSONFormatter{TimestampFormat: DefaultTimestampFormat})

	case PrettyJSONFormat:
		l.Logger.SetFormatter(&log.JSONFormatter{PrettyPrint: true, TimestampFormat: DefaultTimestampFormat})

	default:
		l.Logger.SetFormatter(&log.TextFormatter{TimestampFormat: DefaultTimestampFormat})
	}

	l.logFormat = logFormat
}

// GetFormatter get logger formatter
func GetFormatter() LogFormat { return l.logFormat }

// SetTimestampFormat set timestamp format
func SetTimestampFormat(timestampFormat string) { l.SetTimestampFormat(timestampFormat) }

// SetTimestampFormat set timestamp format
func (l *Logger) SetTimestampFormat(timestampFormat string) {
	switch l.logFormat {
	case JSONFormat:
		l.Logger.SetFormatter(&log.JSONFormatter{TimestampFormat: timestampFormat})

	case PrettyJSONFormat:
		l.Logger.SetFormatter(&log.JSONFormatter{PrettyPrint: true, TimestampFormat: timestampFormat})

	default:
		l.Logger.SetFormatter(&log.TextFormatter{TimestampFormat: timestampFormat})
	}

	l.timestampFormat = timestampFormat
}

// GetTimestampFormat get timestamp format
func GetTimestampFormat() string { return l.timestampFormat }

// SetLogConsole set log console
func SetLogConsole() { l.SetLogConsole() }

// SetLogConsole set log console
func (l *Logger) SetLogConsole() {
	l.SetOutput(os.Stdout)
}

// SetLevel set log level
func SetLevel(level Level) { l.SetLevel(level) }

// SetLevel set log level
func (l *Logger) SetLevel(level Level) {
	l.Logger.SetLevel(log.Level(level))
}

// SetLogger set logger instance
func SetLogger(logger *Logger) { l = logger }

// GetLogger set logger instance
func GetLogger() *Logger { return l }

func getDefaultConfig() *Config {
	return &Config{
		LogLevel:        DebugLevel,
		LogFormat:       JSONFormat,
		TimestampFormat: DefaultTimestampFormat,
	}
}

// NewWithConfig returns a logger instance with custom configuration
func NewWithConfig(config *Config) *Logger {
	if config == nil {
		config = getDefaultConfig()
	}
	logger := log.New()
	l := &Logger{Logger: logger}
	l.logFormat = config.LogFormat
	l.SetFormatter(config.LogFormat)
	if len(config.TimestampFormat) > 0 {
		l.SetTimestampFormat(config.TimestampFormat)
	}

	l.SetLogConsole()
	l.SetLevel(config.LogLevel)

	return l
}

// New returns a logger instance with default configuration
func New() *Logger {
	config := getDefaultConfig()
	logger := log.New()
	l := &Logger{Logger: logger}
	l.SetFormatter(config.LogFormat)
	l.logFormat = config.LogFormat
	if len(config.TimestampFormat) > 0 {
		l.SetTimestampFormat(config.TimestampFormat)
	}

	l.SetLogConsole()
	l.SetLevel(config.LogLevel)

	return l
}

// Trace log with trace level
func Trace(args ...interface{}) {
	l.Trace(args...)
}

// Debug log with debug level
func Debug(args ...interface{}) {
	l.Debug(args...)
}

// Info log with info level
func Info(args ...interface{}) {
	l.Info(args...)
}

// Warn log with warn level
func Warn(args ...interface{}) {
	l.WithField(DefaultSourceField, getCaller()).Warn(args...)
}

// Error log with error level
func Error(args ...interface{}) {
	l.WithField(DefaultSourceField, getCaller()).Error(args...)
}

// Fatal log with fatal level
func Fatal(args ...interface{}) {
	l.WithField(DefaultSourceField, getCaller()).Fatal(args...)
}

// Panic log with panic level
func Panic(args ...interface{}) {
	l.WithField(DefaultSourceField, getCaller()).Panic(args...)
}
