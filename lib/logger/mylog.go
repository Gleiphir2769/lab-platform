package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

type logInfo struct {
	Msg    string
	Fields []zapcore.Field
}

var (
	Warn  chan *logInfo
	Debug chan *logInfo
	Info  chan *logInfo
	Error chan *logInfo
	Fatal chan *logInfo
)

var Log *log
var once sync.Once

type log struct{}

func GetLog() *log {
	once.Do(func() {
		Warn = make(chan *logInfo, LogConfig.BufferSize)
		Debug = make(chan *logInfo, LogConfig.BufferSize)
		Info = make(chan *logInfo, LogConfig.BufferSize)
		Error = make(chan *logInfo, LogConfig.BufferSize)
		Fatal = make(chan *logInfo, LogConfig.BufferSize)
		Log = &log{}
	})
	return Log
}

func (l *log) Fatal(msg string, fields ...Field) {
	Fatal <- &logInfo{
		Msg:    msg,
		Fields: fields,
	}
}

func (l *log) Error(msg string, fields ...Field) {
	Error <- &logInfo{
		Msg:    msg,
		Fields: fields,
	}
}

func (l *log) Info(msg string, fields ...Field) {
	Info <- &logInfo{
		Msg:    msg,
		Fields: fields,
	}
}

func (l *log) Debug(msg string, fields ...Field) {
	Debug <- &logInfo{
		Msg:    msg,
		Fields: fields,
	}
}

func (l *log) Warn(msg string, fields ...Field) {
	Warn <- &logInfo{
		Msg:    msg,
		Fields: fields,
	}
}

func (l *log) Start() {
	go start()
}

func start() {
	for {
		select {
		case lf, _ := <-Warn:
			zap.L().Warn(lf.Msg, lf.Fields...)
		case lf, _ := <-Debug:
			zap.L().Debug(lf.Msg, lf.Fields...)
		case lf, _ := <-Info:
			zap.L().Info(lf.Msg, lf.Fields...)
		case lf, _ := <-Error:
			zap.L().Error(lf.Msg, lf.Fields...)
		case lf, _ := <-Fatal:
			zap.L().Fatal(lf.Msg, lf.Fields...)
		}
	}
}
