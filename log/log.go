// Package log 提供了基于 zap 的高性能日志库
// 支持多级别日志、日志轮转、TraceID 追踪等功能
package log

import (
	"context" // 用于 TraceID 上下文传递
	"fmt"     // 用于字符串格式化
	"os"      // 用于程序退出操作
	"time"    // 用于时间相关操作

	"github.com/daozhonglee/go-util/errorutil"          // 自定义错误处理工具
	rotatelogs "github.com/lestrrat-go/file-rotatelogs" // 日志文件轮转库
	"github.com/petermattis/goid"                       // 获取 goroutine ID
	"go.uber.org/zap"                                   // 高性能日志库
	"go.uber.org/zap/zapcore"                           // zap 核心组件
	"gopkg.in/natefinch/lumberjack.v2"                  // 日志文件切割库
)

const (
	TRACEID = "trace_id" // TraceID 在 context 中的 key，用于请求链路追踪
)

var (
	// Logger 全局日志实例，默认 INFO 级别，输出到 "log" 文件
	Logger = getLogger(zapcore.InfoLevel, "log")
	// atom 原子级别控制器，用于动态调整日志级别
	atom = zap.NewAtomicLevel()
)

// getRotateLogger 创建一个支持文件轮转的日志器
// fmtstr: 日志文件名格式字符串，支持时间模式如 "app.%Y%m%d.log"
// options: rotatelogs 的配置选项，如轮转间隔、保留时间等
func getRotateLogger(fmtstr string, options ...rotatelogs.Option) *zap.SugaredLogger {
	// 创建生产环境的编码器配置
	encoderConf := zap.NewProductionEncoderConfig()
	// 设置时间编码器为空函数（不输出时间）
	encoderConf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
	// 设置级别编码器为空函数（不输出级别）
	encoderConf.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {}
	// 设置调用者编码器为空函数（不输出文件行号）
	encoderConf.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {}
	// 使用控制台编码器创建编码器
	encoder := zapcore.NewConsoleEncoder(encoderConf)

	// 添加默认的最大保留时间：12个月（12*30*24小时）
	options = append(options, rotatelogs.WithMaxAge(12*30*24*time.Hour))
	// 创建轮转日志写入器
	writer, err := rotatelogs.New(fmtstr, options...)
	if err != nil {
		// 如果创建失败，使用自定义错误处理（panic 并打印堆栈）
		errorutil.PanicIfWithStackAndMsg(err, "init logger failed")
	}
	// 创建 zap Core 包装器，配置编码器、写入器和日志级别
	core := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel)
	})

	// 创建生产环境配置
	config := zap.NewProductionConfig()
	// 禁用堆栈跟踪（提高性能）
	config.DisableStacktrace = true
	// 构建日志器
	l, err := config.Build(core)
	if err != nil {
		// 如果构建失败，直接 panic
		panic(fmt.Sprintf("init logger failed err: %v", err))
	}

	// 返回糖化日志器（提供更简单的 API）
	return l.Sugar()
}

// getLogger 创建标准的双文件日志器（DEBUG文件 + WARN文件）
// lv: 日志级别，控制哪些级别的日志会被记录
// path: 日志文件路径前缀，实际会生成 path.DEBUG 和 path.WARN 两个文件
func getLogger(lv zapcore.Level, path string) *zap.SugaredLogger {
	// 获取双核心配置（分别写入不同文件）
	core := getCore(lv, path)
	// 创建生产环境配置
	config := zap.NewProductionConfig()
	// 设置配置的日志级别
	config.Level.SetLevel(lv)
	// 禁用堆栈跟踪（提高性能）
	config.DisableStacktrace = true
	// 构建日志器，跳过一层调用栈（因为有包装函数）
	logger, err := config.Build(core, zap.AddCallerSkip(1))
	if err != nil {
		// 构建失败直接 panic
		panic("init logger failed")
	}
	// 返回糖化日志器
	return logger.Sugar()
}

// SetLevel 动态设置全局日志级别
// level: 新的日志级别，只有大于等于此级别的日志才会被记录
func SetLevel(level zapcore.Level) {
	// 原子性地设置日志级别
	atom.SetLevel(level)
}

// getCore 创建双核心配置，实现日志分级存储
// level: 基础日志级别
// path: 文件路径前缀
func getCore(level zapcore.Level, path string) zap.Option {
	// 创建高优先级日志过滤器，只记录 WARN 及以上级别
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		// 需要同时满足：1) 当前级别允许 2) 是WARN及以上级别
		return level.Enabled(lvl) && lvl >= zapcore.WarnLevel
	})

	// 返回核心包装器，创建双核心 Tee 结构
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(
			// 第一个核心：所有级别的日志写入 .DEBUG 文件
			zapcore.NewCore(getEncoder(), getWriter(fmt.Sprintf("%s.DEBUG", path)), atom),
			// 第二个核心：只有 WARN+ 级别的日志写入 .WARN 文件
			zapcore.NewCore(getEncoder(), getWriter(fmt.Sprintf("%s.WARN", path)), highPriority),
		)
	})
}

// getEncoder 创建自定义的日志编码器，定义日志的输出格式
func getEncoder() zapcore.Encoder {
	// 基于生产环境配置创建编码器配置
	encoderConf := zap.NewProductionEncoderConfig()

	// 自定义时间输出格式：[月-日 时:分:秒.毫秒]
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("01-02 15:04:05.000") + "]")
	}

	// 自定义日志级别显示：[DEBUG/INFO/WARN/ERROR等]
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.String() + "]")
	}

	// 自定义调用者信息输出：[goroutine_id][文件:行号]
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		// 添加 goroutine ID，便于并发调试
		enc.AppendString(fmt.Sprintf("[%d]", goid.Get()))
		// 添加文件路径和行号
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}

	// 应用自定义编码器
	encoderConf.EncodeTime = customTimeEncoder     // 时间格式
	encoderConf.EncodeLevel = customLevelEncoder   // 级别格式
	encoderConf.EncodeCaller = customCallerEncoder // 调用者格式
	encoderConf.ConsoleSeparator = " "             // 字段间分隔符为空格
	// 返回控制台编码器（便于阅读的格式）
	return zapcore.NewConsoleEncoder(encoderConf)
}

// getWriter 创建日志文件写入器，支持文件切割和压缩
// lumberjack这个库目前只支持按文件大小切割(按时间切割效率低且不能保证日志数据不被破坏,详情见https://github.com/natefinch/lumberjack/issues/54）
// 想按日期切割可以使用github.com/lestrrat-go/file-rotatelogs[3]这个库(目前不维护了)

// file: 日志文件路径
func getWriter(file string) zapcore.WriteSyncer {
	// 创建 lumberjack 日志切割器
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,     // 日志文件名
		MaxSize:    1024,     // 文件最大大小（MB），超过后切割
		MaxBackups: 0,        // 保留的备份文件数量，0表示保留所有
		MaxAge:     730 * 10, // 文件最大保留天数（约20年）
		Compress:   true,     // 是否压缩旧的日志文件
	}
	// 将 lumberjack 包装为 zap 的 WriteSyncer
	return zapcore.AddSync(lumberJackLogger)
}

// New 重新初始化全局 Logger
// level: 新的日志级别
// path: 新的日志文件路径前缀
func New(level zapcore.Level, path string) {
	// 重新创建全局 Logger 实例
	Logger = getLogger(level, path)
}

// ========================== 日志函数 ==========================

// CRITICAL 记录严重错误日志但不退出程序
// 内部调用 Panicf 但通过 recover 捕获 panic，同时记录 ERROR 级别日志
// format: 格式化字符串
// a: 格式化参数
func CRITICAL(format string, a ...interface{}) {
	// 使用 defer + recover 捕获 Panicf 产生的 panic
	defer func() {
		if err := recover(); err != nil {
			// 忽略 panic，继续执行
		}
		// 格式化消息并记录为 ERROR 级别
		datastr := fmt.Sprintf(format, a...)
		Logger.Error(datastr)
	}()
	// 调用 Panicf（会产生 panic，但被上面的 recover 捕获）
	Logger.Panicf(format, a...)
}

// Debug 记录 DEBUG 级别日志（调试信息）
// a: 要记录的参数，会被转换为字符串
func Debug(a ...interface{}) {
	Logger.Debug(a...)
}

// Debugf 记录格式化的 DEBUG 级别日志
// format: 格式化字符串
// a: 格式化参数
func Debugf(format string, a ...interface{}) {
	Logger.Debugf(format, a...)
}

// Info 记录 INFO 级别日志（一般信息）
// a: 要记录的参数
func Info(a ...interface{}) {
	Logger.Info(a...)
}

// Infof 记录格式化的 INFO 级别日志
// format: 格式化字符串
// a: 格式化参数
func Infof(format string, a ...interface{}) {
	Logger.Infof(format, a...)
}

// Warn 记录 WARN 级别日志（警告信息）
// a: 要记录的参数
func Warn(a ...interface{}) {
	Logger.Warn(a...)
}

// Warnf 记录格式化的 WARN 级别日志
// format: 格式化字符串
// a: 格式化参数
func Warnf(format string, a ...interface{}) {
	Logger.Warnf(format, a...)
}

// Error 记录 ERROR 级别日志（错误信息）
// a: 要记录的参数
func Error(a ...interface{}) {
	Logger.Error(a...)
}

// Errorf 记录格式化的 ERROR 级别日志
// format: 格式化字符串
// a: 格式化参数
func Errorf(format string, a ...interface{}) {
	Logger.Errorf(format, a...)
}

// Fatal 记录 FATAL 级别日志并退出程序（退出码1）
// a: 要记录的参数
func Fatal(a ...interface{}) {
	Logger.Fatal(a...)
}

// Fatalf 记录格式化的 FATAL 级别日志并退出程序
// format: 格式化字符串
// a: 格式化参数
func Fatalf(format string, a ...interface{}) {
	Logger.Fatalf(format, a...)
}

// FatalfWithExit 记录 FATAL 日志并强制退出程序
// 与 Fatalf 类似，但显式调用 os.Exit(1)
// format: 格式化字符串
// a: 格式化参数
func FatalfWithExit(format string, a ...interface{}) {
	Logger.Fatalf(format, a...) // 记录日志
	os.Exit(1)                  // 强制退出程序
}

// Panic 记录 PANIC 级别日志并引发 panic
// a: 要记录的参数
func Panic(a ...interface{}) {
	Logger.Panic(a...)
}

// Panicf 记录格式化的 PANIC 级别日志并引发 panic
// 注意：会导致程序 panic，除非被 recover 捕获
// format: 格式化字符串
// a: 格式化参数
func Panicf(format string, a ...interface{}) {
	Logger.Panicf(format, a...)
}

// Debugt 记录带 TraceID 的 DEBUG 级别日志，用于请求链路追踪
// ctx: 上下文，如果为 nil 则记录普通 DEBUG 日志
// format: 格式化字符串
// a: 格式化参数
func Debugt(ctx context.Context, format string, a ...interface{}) {
	if ctx == nil {
		// 如果上下文为空，记录普通 DEBUG 日志
		Debugf(format, a...)
		return
	}
	// 在格式化字符串前添加 TraceID 前缀
	format = "[TRACEID:%v] " + format
	// 将 TraceID 插入到参数列表开头
	a = append([]interface{}{ctx.Value(TRACEID)}, a...)
	// 记录带 TraceID 的日志
	Logger.Debugf(format, a...)
}

// Infot 记录带 TraceID 的 INFO 级别日志
// ctx: 上下文，如果为 nil 则记录普通 INFO 日志
// format: 格式化字符串
// a: 格式化参数
func Infot(ctx context.Context, format string, a ...interface{}) {
	if ctx == nil {
		// 如果上下文为空，记录普通 INFO 日志
		Infof(format, a...)
		return
	}
	// 在格式化字符串前添加 TraceID 前缀
	format = "[TRACEID:%v] " + format
	// 将 TraceID 插入到参数列表开头
	a = append([]interface{}{ctx.Value(TRACEID)}, a...)
	// 记录带 TraceID 的日志
	Logger.Infof(format, a...)
}

// Warnt 记录带 TraceID 的 WARN 级别日志
// ctx: 上下文，如果为 nil 则记录普通 WARN 日志
// format: 格式化字符串
// a: 格式化参数
func Warnt(ctx context.Context, format string, a ...interface{}) {
	if ctx == nil {
		// 如果上下文为空，记录普通 WARN 日志
		Warnf(format, a...)
		return
	}
	// 在格式化字符串前添加 TraceID 前缀
	format = "[TRACEID:%v] " + format
	// 将 TraceID 插入到参数列表开头
	a = append([]interface{}{ctx.Value(TRACEID)}, a...)
	// 记录带 TraceID 的日志
	Logger.Warnf(format, a...)
}

// Errort 记录带 TraceID 的 ERROR 级别日志
// ctx: 上下文，如果为 nil 则记录普通 ERROR 日志
// format: 格式化字符串
// a: 格式化参数
func Errort(ctx context.Context, format string, a ...interface{}) {
	if ctx == nil {
		// 如果上下文为空，记录普通 ERROR 日志
		Errorf(format, a...)
		return
	}
	// 在格式化字符串前添加 TraceID 前缀
	format = "[TRACEID:%v] " + format
	// 将 TraceID 插入到参数列表开头
	a = append([]interface{}{ctx.Value(TRACEID)}, a...)
	// 记录带 TraceID 的日志
	Logger.Errorf(format, a...)
}
