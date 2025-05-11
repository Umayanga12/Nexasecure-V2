// Package logger provides a sophisticated logging solution with structured logging,
// multiple log levels, and configurable output formats.
package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for structured logging with various levels
// and contextual logging capabilities.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	WithFields(keysAndValues ...interface{}) Logger
	Sync() error
}

// ZapLogger implements Logger interface using Uber's Zap library.
type ZapLogger struct {
	delegate *zap.SugaredLogger
}

// Config holds logging configuration parameters.
type Config struct {
	Level  string // Log level (debug, info, warn, error, fatal)
	Format string // Log format (json, console)
}

// NewLogger creates a new Logger instance with specified configuration.
func NewLogger(config Config) (Logger, error) {
	logLevel := getZapLevel(config.Level)
	encoderConfig := getEncoderConfig(config.Format)

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(logLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: true,
		Encoding:          config.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	zapLogger, err := zapConfig.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1), // Adjust caller to point to actual log location
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	return &ZapLogger{delegate: zapLogger.Sugar()}, nil
}

// NewConfigFromEnv creates Config from environment variables:
// LOG_LEVEL (default: info), LOG_FORMAT (default: console)
func NewConfigFromEnv() Config {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	format := os.Getenv("LOG_FORMAT")
	if format == "" {
		format = "console"
	}

	return Config{
		Level:  level,
		Format: format,
	}
}

func getZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func getEncoderConfig(format string) zapcore.EncoderConfig {
	if format == "console" {
		cfg := zap.NewDevelopmentEncoderConfig()
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return cfg
	}
	return zap.NewProductionEncoderConfig()
}

// Debug logs a debug message with structured context.
func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.delegate.Debugw(msg, keysAndValues...)
}

// Info logs an info message with structured context.
func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.delegate.Infow(msg, keysAndValues...)
}

// Warn logs a warning message with structured context.
func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.delegate.Warnw(msg, keysAndValues...)
}

// Error logs an error message with structured context.
func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.delegate.Errorw(msg, keysAndValues...)
}

// Fatal logs a fatal message with structured context and exits the program.
func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.delegate.Fatalw(msg, keysAndValues...)
}

// WithFields creates a new logger with additional structured context.
func (l *ZapLogger) WithFields(keysAndValues ...interface{}) Logger {
	return &ZapLogger{delegate: l.delegate.With(keysAndValues...)}
}

// Sync flushes any buffered log entries.
func (l *ZapLogger) Sync() error {
	return l.delegate.Sync()
}