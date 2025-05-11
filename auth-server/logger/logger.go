package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	WithFields(keysAndValues ...interface{}) Logger
	Sync() error
}

type ZapLogger struct {
	delegate *zap.SugaredLogger
}

type Config struct {
	Level  string 
	Format string 
}

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
		zap.AddCallerSkip(1), 
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	return &ZapLogger{delegate: zapLogger.Sugar()}, nil
}

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

func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.delegate.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.delegate.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.delegate.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.delegate.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.delegate.Fatalw(msg, keysAndValues...)
}

func (l *ZapLogger) WithFields(keysAndValues ...interface{}) Logger {
	return &ZapLogger{delegate: l.delegate.With(keysAndValues...)}
}

func (l *ZapLogger) Sync() error {
	return l.delegate.Sync()
}