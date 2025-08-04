package log

import "strings"

// 全局单例实例
var globalLoggerIFace = &LoggerIFace{}

func ILogger() *LoggerIFace {
	return globalLoggerIFace
}

type LoggerIFace struct{}

func (*LoggerIFace) Print(v ...interface{}) {
	Info(v...)
}

func (*LoggerIFace) Println(v ...interface{}) {
	Info(v...)
}

func (*LoggerIFace) Printf(format string, v ...interface{}) {
	format = strings.TrimRight(format, "\n")
	Infof(format, v...)
}
