package log

import (
	"testing"

	"go.uber.org/zap"
)

// JSON输出到STDERR，Infovel及以上，包括时间戳和caller信息
func TestZapDevelopment(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	logger.Info("info")
	logger.Debug("debug")
	logger.Warn("warn")
	logger.Error("error")
	logger.Fatal("fatal")
	logger.Panic("panic")
}

// 控制台输出到stderr，debuglevel及以上，包括时间戳和caller信息
func TestZapProduction(t *testing.T) {
	logger, _ := zap.NewProduction()
	logger.Info("info")
	logger.Debug("debug")
	logger.Warn("warn")
	logger.Error("error")
	logger.Fatal("fatal")
	logger.Panic("panic")
}

// JSON输出到Stdout，可用于示例和测试
func TestZapExample(t *testing.T) {
	logger := zap.NewExample()
	logger.Info("info")
	logger.Debug("debug")
	logger.Warn("warn")
	logger.Error("error")
	logger.Fatal("fatal")
	logger.Panic("panic")
}
