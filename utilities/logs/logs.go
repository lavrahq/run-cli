package logs

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log provides standard logging functionality.
var Log *zap.Logger

// NewLogger instantiates a new logger instance.
func NewLogger() (*zap.Logger, error) {
	var LogPath, _ = homedir.Expand("~/.lavra/cli.log")

	cfg := zap.NewProductionConfig()
	cfg.Development = !viper.GetBool("mode.production")
	cfg.OutputPaths = []string{LogPath, "stdout"}
	cfg.Encoding = "console"
	cfg.EncoderConfig.LineEnding = "\r\n"

	if viper.GetBool("mode.debug") {
		cfg.Level.SetLevel(zapcore.DebugLevel)
	} else {
		if viper.GetBool("mode.verbose") {
			cfg.Level.SetLevel(zapcore.InfoLevel)
		} else {
			cfg.Level.SetLevel(zapcore.WarnLevel)
		}
	}

	return cfg.Build()
}

// InitGlobalLogging initializes the global Log var.
func InitGlobalLogging() *zap.Logger {
	Log, _ = NewLogger()

	return Log
}
