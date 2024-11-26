package logger

import (
	"OnlineMusic/utils"
	"io"
	"log"
	"log/slog"
	"os"
)

type Logger = slog.Logger

func New(logLevelFlag, logFilePath string) *Logger {
	logLevel := parseLogLevel(logLevelFlag)

	options := &slog.HandlerOptions{
		Level: logLevel,
	}

	if logLevel == slog.LevelDebug {
		options.AddSource = true
	}

	file, err := utils.FileHandle(logFilePath)
	if err != nil {
		log.Fatalf("failed to open: %v", err)
	}

	w := io.MultiWriter(file, os.Stdout)
	logger := slog.New(slog.NewJSONHandler(w, options))
	slog.SetDefault(logger)

	return logger
}

func parseLogLevel(logLevelFlag string) slog.Level {
	logLevelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	logLevel, ok := logLevelMap[logLevelFlag]

	if !ok {
		logLevel = slog.LevelInfo
	}

	return logLevel
}
