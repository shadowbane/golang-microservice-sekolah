package logger

import (
	"flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"strconv"
)

type LogConfig struct {
	// Level is the log level to log at
	Level string

	// ConsoleEnabled enable console logging
	ConsoleEnabled bool

	// ConsoleJson is a flag to output log messages in json format
	ConsoleJson bool

	// FileEnabled enable file logging
	FileEnabled bool

	// FileJson is a flag to output log messages in json format
	FileJson bool

	// Directory to log to when filelogging is enabled
	Directory string

	// Filename is the name of the logfile which will be placed inside the directory
	Filename string

	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int

	// MaxBackups the max number of rolled files to keep
	MaxBackups int

	// MaxAge the max age in days to keep a logfile
	MaxAge int
}

func LoadEnvForLogger() *LogConfig {
	lc := &LogConfig{}

	ConsoleEnabled, _ := strconv.ParseBool(getenv("LOG_CONSOLE_ENABLED", "true"))
	ConsoleJson, _ := strconv.ParseBool(getenv("LOG_CONSOLE_JSON", "false"))
	FileEnabled, _ := strconv.ParseBool(getenv("LOG_FILE_ENABLED", "false"))
	FileJson, _ := strconv.ParseBool(getenv("LOG_FILE_JSON", "false"))
	MaxSize, _ := strconv.Atoi(getenv("LOG_MAX_SIZE", "25"))
	MaxBackups, _ := strconv.Atoi(getenv("LOG_MAX_BACKUPS", "5"))
	MaxAge, _ := strconv.Atoi(getenv("LOG_MAX_AGE", "30"))

	flag.StringVar(&lc.Level, "Level", getenv("LOG_LEVEL", "info"), "Log level to log at")
	flag.BoolVar(&lc.ConsoleEnabled, "ConsoleEnabled", ConsoleEnabled, "Enable console logging")
	flag.BoolVar(&lc.ConsoleJson, "ConsoleJson", ConsoleJson, "Output log messages in json format")
	flag.BoolVar(&lc.FileEnabled, "FileEnabled", FileEnabled, "Enable file logging")
	flag.BoolVar(&lc.FileJson, "FileJson", FileJson, "Output log messages in json format")
	flag.StringVar(&lc.Directory, "Directory", getenv("LOG_DIRECTORY", "log"), "Directory to log to when filelogging is enabled")
	flag.StringVar(&lc.Filename, "Filename", getenv("LOG_FILENAME", "app.log"), "Name of the logfile which will be placed inside the directory")
	flag.IntVar(&lc.MaxSize, "MaxSize", MaxSize, "Max size in MB of the logfile before it's rolled")
	flag.IntVar(&lc.MaxBackups, "MaxBackups", MaxBackups, "Max number of rolled files to keep")
	flag.IntVar(&lc.MaxAge, "MaxAge", MaxAge, "Max age in days to keep a logfile")

	return lc
}

// getenv get environment variable or fallback to default value if not set
func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func (c *LogConfig) ConsoleEncoderConfig() zapcore.EncoderConfig {
	encoderType := zapcore.CapitalColorLevelEncoder

	if c.ConsoleJson {
		encoderType = zapcore.CapitalLevelEncoder
	}

	return zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    encoderType,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

func (c *LogConfig) FileEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// InitLumberjackLogger initialize lumberjack logger.
// Automatically create log directory if not exists.
// This will roll the log file when it reaches the max size.
func (c *LogConfig) InitLumberjackLogger() lumberjack.Logger {
	if err := os.MkdirAll(c.Directory, 0744); err != nil {
		zap.S().Fatalf("can't create log directory %s: %s", c.Directory, err)
	}

	return lumberjack.Logger{
		Filename:   path.Join(c.Directory, c.Filename),
		MaxSize:    c.MaxSize, //MB
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge, //days
		Compress:   false,
	}
}
