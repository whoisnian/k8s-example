package global

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LOG *zap.Logger

func SetupLogger() {
	var err error
	if CFG.Debug {
		LOG, err = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
			Development: true,
			Encoding:    "console",
			EncoderConfig: zapcore.EncoderConfig{ // from zap.NewDevelopmentEncoderConfig()
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "M",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			DisableStacktrace: true,
			OutputPaths:       []string{"stderr"},
			ErrorOutputPaths:  []string{"stderr"},
		}.Build()
	} else {
		LOG, err = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Encoding:    "json",
			EncoderConfig: zapcore.EncoderConfig{ // from zap.NewProductionEncoderConfig()
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.EpochTimeEncoder,
				EncodeDuration: zapcore.MillisDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			DisableStacktrace: true,
			OutputPaths:       []string{"stderr"},
			ErrorOutputPaths:  []string{"stderr"},
		}.Build()
	}
	if err != nil {
		panic(err)
	}
}
