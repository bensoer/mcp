package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	PRODUCTION  = "production"
	DEVELOPMENT = "development"
)

// InitLogger initializes a global Zap logger.
// Set env to "production" for JSON logs or "development" for console logs.
func InitLogger(env string) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human-readable timestamps
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Colored log levels
	}

	logger, err := config.Build(zap.AddCallerSkip(1)) // AddCallerSkip ensures correct file line tracking
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}

	// Replace global instances so zap.L() and zap.S() use this configuration
	zap.ReplaceGlobals(logger)
}
