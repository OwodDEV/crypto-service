package logger

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/OwodDEV/crypto-service/internal/config"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(cfg *config.Config) (logger *Logger) {
	logLevel := &slog.LevelVar{}
	switch cfg.Logger.Level {
	case "DEBUG":
		logLevel.Set(slog.LevelDebug)
	case "INFO":
		logLevel.Set(slog.LevelInfo)
	case "WARN":
		logLevel.Set(slog.LevelWarn)
	case "ERROR":
		logLevel.Set(slog.LevelError)
	default:
		log.Fatal("invalid log level: " + cfg.Logger.Level)
	}

	// Configure Lumberjack
	fileLogger := &lumberjack.Logger{
		Filename:   cfg.Logger.Filename,
		MaxSize:    cfg.Logger.MaxSize,
		MaxBackups: cfg.Logger.MaxBackups,
		MaxAge:     cfg.Logger.MaxAge,
		Compress:   cfg.Logger.Compress,
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	multiWriter := io.MultiWriter(os.Stdout, fileLogger)
	baseLogger := slog.New(slog.NewJSONHandler(multiWriter, opts))

	logger = &Logger{baseLogger}
	slog.SetDefault(baseLogger)

	return
}
