package logging

import (
	"api/constant"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// Init sets the default slog logger for the process and returns it.
//
// Env/config (via viper):
// - LOG_LEVEL: debug|info|warn|error (default: info)
// - LOG_ADD_SOURCE: true|false (default: false)
func Init() *slog.Logger {
	level := parseLevel()
	addSource := viper.GetBool("LOG_ADD_SOURCE")

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func parseLevel() slog.Level {
	switch constant.GetEnvironment() {
	case constant.EnvironmentLocal, constant.EnvironmentDevelopment:
		return slog.LevelDebug
	case constant.EnvironmentStaging:
		return slog.LevelWarn
	case constant.EnvironmentProduction:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
