package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lg *zap.Logger

// Init
func Init() (err error) {
	writeSyncer := getLogWriter(LogConfig.Filename, LogConfig.MaxSize, LogConfig.MaxBackups, LogConfig.MaxAge)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(LogConfig.Level))
	if err != nil {
		return
	}

	var allCore []zapcore.Core
	if LogConfig.Level == "debug" || LogConfig.Level == "info" {
		allCore = append(allCore, zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), l))
	}
	allCore = append(allCore, zapcore.NewCore(encoder, writeSyncer, l))

	core := zapcore.NewTee(allCore...)

	lg = zap.New(core, zap.AddCaller())
	// Replace the global logger instance in the ZAP package, and then use only `zap.L()` calls in other packages
	zap.ReplaceGlobals(lg)
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
