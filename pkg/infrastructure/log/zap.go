package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLogConfig = "debug"
)

func NewLogger(envLogConfig string) (*zap.Logger, error) {

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.TimeKey = "time"

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:      "json",
		EncoderConfig: encoderCfg,
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	if envLogConfig == DebugLogConfig {
		cfg.Development = true
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	// https://github.com/uber-go/zap/issues/328
	defer logger.Sync()

	return logger, nil
}
