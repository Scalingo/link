package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Default generate a logrus logger with the configuration defined in the environment and the hooks used in the plugins
func Default() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logLevel())
	logger.Formatter = formatter()

	for _, hook := range Plugins().Hooks() {
		logger.Hooks.Add(hook)
	}

	return logger
}

// NewContextWithLogger generate a new context (based on context.Background()) and add a Default() logger on top of it
func NewContextWithLogger() context.Context {
	return AddLoggerToContext(context.Background())
}

// AddLoggerToContext add the Default() logger on top of the current context
func AddLoggerToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "logger", logrus.NewEntry(Default()))
}

// Get return the logger stored in the context or create a new one if the logger is not set
func Get(ctx context.Context) logrus.FieldLogger {
	if logger, ok := ctx.Value("logger").(logrus.FieldLogger); ok {
		return logger
	}

	return Default().WithField("invalid_context", true)
}

// ToCtx add a logger to a context
func ToCtx(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}
