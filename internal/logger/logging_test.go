package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitLogger_Development(t *testing.T) {
	InitLogger(DEVELOPMENT)

	if zap.L() == nil {
		t.Fatal("zap.L() is nil after InitLogger")
	}

	if !zap.L().Core().Enabled(zapcore.DebugLevel) {
		t.Error("development logger should have DebugLevel enabled")
	}
}

func TestInitLogger_Production(t *testing.T) {
	InitLogger(PRODUCTION)

	if zap.L() == nil {
		t.Fatal("zap.L() is nil after InitLogger")
	}

	if zap.L().Core().Enabled(zapcore.DebugLevel) {
		t.Error("production logger should not have DebugLevel enabled")
	}
}

func TestInitLogger_SugaredLogger(t *testing.T) {
	InitLogger(DEVELOPMENT)

	if zap.S() == nil {
		t.Fatal("zap.S() is nil after InitLogger")
	}

	zap.S().Info("test sugared logger message")
}
