package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger is the global logger instance
	Logger *zap.Logger
)

// InitLogger initializes the global logger instance
func InitLogger(logLevel string, logFormat string) error {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return err
	}

	var encoderConfig zapcore.EncoderConfig
	if logFormat == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.Lock(zapcore.AddSync(zapcore.AddSync(os.Stdout))),
		level,
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	return nil
}