package log

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type SugoWriter struct {
	mlog *zap.SugaredLogger
}

// 实现gorm/logger.Writer接口
func (m *SugoWriter) Printf(format string, v ...interface{}) {
	format = strings.Replace(format, "\n", "", -1)
	logstr := fmt.Sprintf(format, v...)
	Debug(logstr)
}

func NewSugoWriter() *SugoWriter {
	return &SugoWriter{mlog: Logger}
}
