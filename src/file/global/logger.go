package global

import (
	"go.uber.org/zap"
)

var LOG *zap.Logger

func SetupLogger() {
	var err error
	if CFG.Debug {
		LOG, err = zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}.Build()
	} else {
		LOG, err = zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}.Build()
	}
	if err != nil {
		panic(err)
	}
}
