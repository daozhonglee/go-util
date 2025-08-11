package log

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zap 的一些基础API 测试
var sugarLogger *zap.SugaredLogger

func TestZap(t *testing.T) {
	InitLogger()
	defer sugarLogger.Sync()
	sugarLogger.Info("info")
	sugarLogger.Debug("debug")
	sugarLogger.Warn("warn")
	sugarLogger.Error("error")
	sugarLogger.Fatal("fatal")
	sugarLogger.Panic("panic")
}
func InitLogger() {
	sugarLogger = getLogger(zapcore.DebugLevel, "./test.log")
}

func getEncoder2() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	//如果想要追加写入可以查看我的博客文件操作那一章
	file, _ := os.Create("./test.log")
	return zapcore.AddSync(file)
}

func TestNewProductionLogger(t *testing.T) {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.DisableCaller = true
	config.Encoding = "json"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	logger, _ := config.Build(zap.AddCallerSkip(1))
	logger.Sugar().Info("info")
	logger.Sugar().Debug("debug")
	logger.Sugar().Warn("warn")
	logger.Sugar().Error("error")
	logger.Sugar().Fatal("fatal")
	logger.Sugar().Panic("panic")
}
func TestNewDevelopmentLogger(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	logger.Sugar().Info("info")
	logger.Sugar().Debug("debug")
	logger.Sugar().Warn("warn")
	logger.Sugar().Error("error")
}
