package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/config"
	"github.com/rs/zerolog"
)

var (
	cfg = config.ConfigSingleton.GetInstance()
)

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	once     sync.Once
	instance *Logger
	logger   *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

func NewLogger(level string) *Logger {
	logLevel, found := map[string]zerolog.Level{
		"error": zerolog.ErrorLevel,
		"warn":  zerolog.WarnLevel,
		"info":  zerolog.InfoLevel,
		"debug": zerolog.DebugLevel,
	}[strings.ToLower(level)]

	if !found {
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stdout).
		Level(logLevel).
		With().
		Timestamp()

	if logLevel == zerolog.DebugLevel {
		logger = logger.CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 3)
	}

	log := logger.Logger()

	return &Logger{logger: &log}
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.formatAndLog("info", message, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.formatAndLog("warning", message, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

func (l *Logger) Error(message interface{}, args ...interface{}) {
	if l.logger.GetLevel() == zerolog.DebugLevel {
		l.Debug(message, args...)
	}
	l.msg("error", message, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)
	os.Exit(1)
}

func (l *Logger) formatAndLog(level string, message string, args ...interface{}) {
	switch level {
	case "info":
		if len(args) == 0 {
			l.logger.Info().Msg(message)
		} else {
			l.logger.Info().Msgf(message, args...)
		}
	case "debug":
		if len(args) == 0 {
			l.logger.Debug().Msg(message)
		} else {
			l.logger.Debug().Msgf(message, args...)
		}
	case "warn":
		if len(args) == 0 {
			l.logger.Warn().Msg(message)
		} else {
			l.logger.Warn().Msgf(message, args...)
		}
	case "error":
		if len(args) == 0 {
			l.logger.Error().Msg(message)
		} else {
			l.logger.Error().Msgf(message, args...)
		}
	case "fatal":
		if len(args) == 0 {
			l.logger.Fatal().Msg(message)
		} else {
			l.logger.Fatal().Msgf(message, args...)
		}
	}
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch v := message.(type) {
	case error:
		l.formatAndLog(level, v.Error(), args...)
	case string:
		l.formatAndLog(level, v, args...)
	default:
		l.formatAndLog(level, fmt.Sprintf("%s message %v with unknown type %T", level, message, v), args...)
	}
}

func (log *Logger) GetInstance() *Logger {
	log.once.Do(func() {
		log.instance = NewLogger(cfg.Log.Level)
	})
	return log.instance
}

var LoggerSingleton Logger
