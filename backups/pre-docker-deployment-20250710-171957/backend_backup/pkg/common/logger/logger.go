package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents logging levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns string representation of log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a logger instance
type Logger struct {
	level  LogLevel
	logger *log.Logger
	fields map[string]interface{}
}

// New creates a new logger instance
func New(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
		fields: make(map[string]interface{}),
	}
}

// NewWithOutput creates a new logger with custom output
func NewWithOutput(level LogLevel, output *os.File) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(output, "", 0),
		fields: make(map[string]interface{}),
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	if l.level <= DEBUG {
		l.log(DEBUG, msg)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	if l.level <= INFO {
		l.log(INFO, msg)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	if l.level <= WARN {
		l.log(WARN, msg)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	if l.level <= ERROR {
		l.log(ERROR, msg)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	l.log(FATAL, msg)
	os.Exit(1)
}

// WithField adds a field to the logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &Logger{
		level:  l.level,
		logger: l.logger,
		fields: newFields,
	}
}

// WithFields adds multiple fields to the logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &Logger{
		level:  l.level,
		logger: l.logger,
		fields: newFields,
	}
}

// WithError adds an error to the logger context
func (l *Logger) WithError(err error) *Logger {
	return l.WithField("error", err.Error())
}

// log performs the actual logging
func (l *Logger) log(level LogLevel, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	// Build the log message
	logMsg := fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), msg)
	
	// Add fields if any
	if len(l.fields) > 0 {
		fieldsStr := ""
		for k, v := range l.fields {
			if fieldsStr != "" {
				fieldsStr += ", "
			}
			fieldsStr += fmt.Sprintf("%s=%v", k, v)
		}
		logMsg += fmt.Sprintf(" | %s", fieldsStr)
	}
	
	l.logger.Println(logMsg)
}

// ParseLogLevel parses a string log level
func ParseLogLevel(level string) LogLevel {
	switch level {
	case "debug", "DEBUG":
		return DEBUG
	case "info", "INFO":
		return INFO
	case "warn", "WARN", "warning", "WARNING":
		return WARN
	case "error", "ERROR":
		return ERROR
	case "fatal", "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// Global logger instance
var defaultLogger = New(INFO)

// Global logging functions
func Debug(msg string) {
	defaultLogger.Debug(msg)
}

func Info(msg string) {
	defaultLogger.Info(msg)
}

func Warn(msg string) {
	defaultLogger.Warn(msg)
}

func Error(msg string) {
	defaultLogger.Error(msg)
}

func Fatal(msg string) {
	defaultLogger.Fatal(msg)
}

func WithField(key string, value interface{}) *Logger {
	return defaultLogger.WithField(key, value)
}

func WithFields(fields map[string]interface{}) *Logger {
	return defaultLogger.WithFields(fields)
}

func WithError(err error) *Logger {
	return defaultLogger.WithError(err)
}

func SetLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	defaultLogger = logger
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	return defaultLogger
}
