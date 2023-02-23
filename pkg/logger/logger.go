package logger

import (
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	FileDir    string `yaml:"FileDir"`
	MaxSize    int    `yaml:"MaxSize"`
	MaxBackups int    `yaml:"MaxBackups"`
	MaxAge     int    `yaml:"MaxAge"`
	Compress   bool   `yaml:"Compress"`
	LocalTime  bool   `yaml:"LocalTime"`

	Level zapcore.Level
}

func (c *Config) SetLevel(level zapcore.Level) *Config {
	c.Level = level
	return c
}

func NewLogger(c *Config) *zap.SugaredLogger {
	return newLogger(c.SetLevel(zap.InfoLevel))
}

func newLogger(c *Config) *zap.SugaredLogger {
	writeSyncer := getLogWriter(c)
	encoder := getEncoder()
	fileCore := zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)
	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
	core := zapcore.NewTee(fileCore, consoleCore)
	return zap.New(core).Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(c *Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(c.FileDir, "current.log"),
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}
