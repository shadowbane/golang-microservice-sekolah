package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
)

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

// initEncoder initializes the encoder for the console and file logging
func initEncoder(c *LogConfig) (zapcore.Encoder, zapcore.Encoder) {
	fileEncoderConfig := c.FileEncoderConfig()
	consoleEncoderConfig := c.ConsoleEncoderConfig()

	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig)

	if c.ConsoleJson {
		log.Print("JSON output enabled for ConsoleJson")
		consoleEncoder = zapcore.NewJSONEncoder(consoleEncoderConfig)
	}
	if c.FileJson {
		log.Print("JSON output enabled for FileJson")
		fileEncoder = zapcore.NewJSONEncoder(fileEncoderConfig)
	}

	return consoleEncoder, fileEncoder
}

// initZapLog initializes the zap logger
func initZapLog(logLevel zapcore.Level, c *LogConfig) *zap.Logger {
	ll := c.InitLumberjackLogger()

	consoleEncoder, fileEncoder := initEncoder(c)

	var cores []zapcore.Core

	if c.ConsoleEnabled {
		log.Print("Console logging is enabled")
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), logLevel))
	}
	if c.FileEnabled {
		log.Print("File logging is enabled")
		cores = append(cores, zapcore.NewCore(fileEncoder, lumberjackSink{Logger: &ll}, logLevel))
	}
	core := zapcore.NewTee(cores...)

	unsugared := zap.New(core)

	return unsugared
}

func Init(c *LogConfig) {
	zap.S().Debug("initializing logger")
	logLevel := getLogLevel(c.Level)
	logManager := initZapLog(logLevel, c)

	zap.ReplaceGlobals(logManager)

	defer func(logManager *zap.Logger) {
		err := logManager.Sync()
		if err != nil {

		}
	}(logManager) // flushes buffer, if any
}

// getLogLevel returns the log level based on the config.
//
// Valid log levels are debug, info, warn, and error.
// If an invalid log level is provided, info is returned
func getLogLevel(configLevel string) zapcore.Level {
	switch configLevel {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
