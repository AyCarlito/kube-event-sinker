package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ContextKey string

const Logger ContextKey = "logger"

// ContextWithLogger returns a copy of the provided context.Context that associates a key with the provided *zap.Logger.
func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, Logger, logger)
}

// LoggerFromContext returns the *zap.Logger associated with a key in the provided context.Context.
func LoggerFromContext(ctx context.Context) *zap.Logger {
	return ctx.Value(Logger).(*zap.Logger)
}

// NewZapConfig returns a new zap.Config for constructing a zap.Logger.
// The default production configuration serves as a base, with overrides on specific fields.
func NewZapConfig() zap.Config {
	z := zap.NewProductionConfig()
	z.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	z.EncoderConfig.TimeKey = "time"
	return z
}
