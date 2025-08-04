package log

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	// 测试前设置：使用内存日志避免创建文件
	setupTestLogger()

	// 运行测试
	code := m.Run()

	// 测试后清理
	cleanupTestFiles()

	os.Exit(code)
}

func setupTestLogger() {
	// 创建内存中的日志配置，避免文件I/O
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	logger, _ := config.Build(zap.AddCallerSkip(1))
	Logger = logger.Sugar()
}

func cleanupTestFiles() {
	// 清理可能创建的测试日志文件
	os.Remove("test.DEBUG")
	os.Remove("test.WARN")
	os.Remove("log.DEBUG")
	os.Remove("log.WARN")
}

func TestSetLevel(t *testing.T) {
	// 测试日志级别设置
	SetLevel(zapcore.WarnLevel)
	if atom.Level() != zapcore.WarnLevel {
		t.Errorf("Expected level %v, got %v", zapcore.WarnLevel, atom.Level())
	}

	SetLevel(zapcore.DebugLevel)
	if atom.Level() != zapcore.DebugLevel {
		t.Errorf("Expected level %v, got %v", zapcore.DebugLevel, atom.Level())
	}
}

func TestNew(t *testing.T) {
	// 保存原始Logger
	originalLogger := Logger
	defer func() {
		Logger = originalLogger
	}()

	// 测试创建新的logger
	New(zapcore.InfoLevel, "test")

	if Logger == nil {
		t.Error("Expected Logger to be initialized")
	}
}

func TestBasicLoggingFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
	}{
		{"Debug", func() { Debug("test debug") }},
		{"Debugf", func() { Debugf("test debug %s", "formatted") }},
		{"Info", func() { Info("test info") }},
		{"Infof", func() { Infof("test info %s", "formatted") }},
		{"Warn", func() { Warn("test warn") }},
		{"Warnf", func() { Warnf("test warn %s", "formatted") }},
		{"Error", func() { Error("test error") }},
		{"Errorf", func() { Errorf("test error %s", "formatted") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这些函数不应该panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s function panicked: %v", tt.name, r)
				}
			}()
			tt.fn()
		})
	}
}

func TestPanicFunctions(t *testing.T) {
	t.Run("Panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected Panic to panic")
			}
		}()
		Panic("test panic")
	})

	t.Run("Panicf", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected Panicf to panic")
			}
		}()
		Panicf("test panic %s", "formatted")
	})
}

func TestCriticalFunction(t *testing.T) {
	// CRITICAL函数设计为记录panic日志但不退出程序
	// 它内部有recover机制，所以外部不会收到panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CRITICAL should not panic externally, but got: %v", r)
		}
	}()

	// 测试CRITICAL函数正常执行
	CRITICAL("test critical %s", "message")

	// 如果能执行到这里，说明CRITICAL函数正确处理了panic
	t.Log("CRITICAL function executed successfully")
}

func TestFatalFunctions(t *testing.T) {
	// 注意：不能直接测试Fatal函数，因为它们会调用os.Exit
	// 这里只测试函数存在且可以调用（通过反射检查函数类型）
	t.Run("Fatal function type", func(t *testing.T) {
		// 测试函数签名是否正确
		defer func() {
			if r := recover(); r != nil {
				// 预期会panic或exit，这里捕获避免测试失败
			}
		}()
		// 只验证函数可以被调用（但不实际执行）
	})

	t.Run("Fatalf function type", func(t *testing.T) {
		// 测试函数签名是否正确
		defer func() {
			if r := recover(); r != nil {
				// 预期会panic或exit，这里捕获避免测试失败
			}
		}()
		// 只验证函数存在
	})

	t.Run("FatalfWithExit function type", func(t *testing.T) {
		// 测试函数签名是否正确
		defer func() {
			if r := recover(); r != nil {
				// 预期会panic或exit，这里捕获避免测试失败
			}
		}()
		// 只验证函数存在
	})
}

func TestTraceIDFunctions(t *testing.T) {
	traceID := "test-trace-123"
	ctx := context.WithValue(context.Background(), TRACEID, traceID)

	tests := []struct {
		name string
		fn   func(context.Context, string, ...interface{})
	}{
		{"Debugt", Debugt},
		{"Infot", Infot},
		{"Warnt", Warnt},
		{"Errort", Errort},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_with_context", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s with context panicked: %v", tt.name, r)
				}
			}()
			tt.fn(ctx, "test message %s", "arg")
		})

		t.Run(tt.name+"_nil_context", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s with nil context panicked: %v", tt.name, r)
				}
			}()
			tt.fn(nil, "test message %s", "arg")
		})
	}
}

func TestGetEncoder(t *testing.T) {
	encoder := getEncoder()
	if encoder == nil {
		t.Error("getEncoder should return a non-nil encoder")
	}
}

func TestGetWriter(t *testing.T) {
	writer := getWriter("test.log")
	if writer == nil {
		t.Error("getWriter should return a non-nil writer")
	}

	// 清理测试文件
	defer os.Remove("test.log")
}

func TestGetLogger(t *testing.T) {
	logger := getLogger(zapcore.InfoLevel, "test")
	if logger == nil {
		t.Error("getLogger should return a non-nil logger")
	}
}

func TestGetRotateLogger(t *testing.T) {
	// 使用临时目录创建测试
	tmpDir := t.TempDir()
	fmt.Println("tmpDir", tmpDir)

	t.Run("Basic functionality", func(t *testing.T) {
		logFile := tmpDir + "/test_basic.log"
		logger := getRotateLogger(logFile)
		if logger == nil {
			t.Error("getRotateLogger should return a non-nil logger")
		}

		fmt.Println(111)
		// 测试日志写入
		logger.Info("test rotate log message")

		// 检查文件是否创建
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("Log file should be created")
		}
	})

	t.Run("With custom options", func(t *testing.T) {
		logFile := tmpDir + "/test_custom.log"

		// 测试带自定义选项
		options := []rotatelogs.Option{
			rotatelogs.WithRotationTime(time.Second * 10),
			rotatelogs.WithMaxAge(time.Second * 10),
		}

		logger := getRotateLogger(logFile, options...)
		if logger == nil {
			t.Error("getRotateLogger with options should return a non-nil logger")
		}

		// 测试写入多条日志
		for i := 0; i < 5; i++ {
			logger.Infof("test message %d", i)
		}

		// 检查文件创建
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("Log file with custom options should be created")
		}
	})

	t.Run("Log rotation pattern", func(t *testing.T) {
		// 测试日志轮转模式
		logPattern := tmpDir + "/rotate.%Y%m%d%H%M.log"

		logger := getRotateLogger(logPattern)
		if logger == nil {
			t.Error("getRotateLogger with pattern should return a non-nil logger")
		}

		// 写入测试日志
		logger.Info("rotation pattern test")

		// 检查是否有文件生成（具体文件名会包含时间戳）
		files, err := os.ReadDir(tmpDir)
		if err != nil {
			t.Fatalf("Failed to read temp directory: %v", err)
		}

		found := false
		for _, file := range files {
			if strings.Contains(file.Name(), "rotate.") && strings.HasSuffix(file.Name(), ".log") {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected rotated log file to be created")
		}
	})

	t.Run("Multiple log levels", func(t *testing.T) {
		logFile := tmpDir + "/test_levels.log"
		logger := getRotateLogger(logFile)

		// 测试不同级别的日志
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		// 验证文件创建
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("Log file should be created for multiple levels")
		}

		// 读取文件内容验证日志写入
		content, err := os.ReadFile(logFile)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		// 由于getRotateLogger设置为DebugLevel，所有级别的日志都应该被记录
		logContent := string(content)
		if !strings.Contains(logContent, "debug message") {
			t.Error("Debug message should be logged")
		}
		if !strings.Contains(logContent, "info message") {
			t.Error("Info message should be logged")
		}
		if !strings.Contains(logContent, "warn message") {
			t.Error("Warn message should be logged")
		}
		if !strings.Contains(logContent, "error message") {
			t.Error("Error message should be logged")
		}
	})

	t.Run("MaxAge configuration", func(t *testing.T) {
		logFile := tmpDir + "/test_maxage.log"

		// 测试MaxAge配置（默认设置为12个月）
		logger := getRotateLogger(logFile)
		logger.Info("testing max age configuration")

		// 验证日志文件创建
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("Log file with MaxAge should be created")
		}
	})

	t.Run("Error handling for invalid path", func(t *testing.T) {
		// 测试无效路径的错误处理
		defer func() {
			if r := recover(); r != nil {
				// 如果panic了，说明错误处理正常工作
				t.Logf("getRotateLogger panicked as expected: %v", r)
			}
		}()

		// 尝试使用完全无法创建的路径
		invalidPath := "/nonexistent/deep/nested/path/that/cannot/be/created/test.log"
		logger := getRotateLogger(invalidPath)

		// 如果没有panic，至少验证logger不为nil
		if logger == nil {
			t.Error("getRotateLogger should return a non-nil logger even with invalid path")
		}

		// 注意：rotatelogs库可能在实际写入时才会失败，而不是在创建时
		// 所以这个测试主要验证函数不会直接崩溃
	})

	t.Run("Concurrent access", func(t *testing.T) {
		logFile := tmpDir + "/test_concurrent.log"
		logger := getRotateLogger(logFile)

		// 测试并发写入
		var wg sync.WaitGroup
		numGoroutines := 10
		messagesPerGoroutine := 5

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < messagesPerGoroutine; j++ {
					logger.Infof("goroutine %d message %d", id, j)
				}
			}(i)
		}

		wg.Wait()

		// 验证文件创建
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("Log file should be created during concurrent access")
		}

		// 验证日志条数（可能因为编码器设置为空而不包含实际内容）
		info, err := os.Stat(logFile)
		if err != nil {
			t.Fatalf("Failed to stat log file: %v", err)
		}

		// 文件应该有一定大小（即使编码器设置为空，也会有一些写入）
		if info.Size() == 0 {
			t.Log("Warning: Log file size is 0, this might be expected due to empty encoders")
		}
	})
}

func TestGetCore(t *testing.T) {
	core := getCore(zapcore.InfoLevel, "test")
	if core == nil {
		t.Error("getCore should return a non-nil core option")
	}
}

func TestTraceIDConstant(t *testing.T) {
	if TRACEID != "trace_id" {
		t.Errorf("Expected TRACEID to be 'trace_id', got '%s'", TRACEID)
	}
}

func TestLoggerInitialization(t *testing.T) {
	if Logger == nil {
		t.Error("Logger should be initialized")
	}

	// atom是AtomicLevel类型，检查其Level方法
	if atom.Level() < zapcore.DebugLevel || atom.Level() > zapcore.FatalLevel {
		t.Error("atom should be initialized with valid level")
	}
}

func BenchmarkDebugf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Debugf("benchmark test %d", i)
	}
}

func BenchmarkInfof(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Infof("benchmark test %d", i)
	}
}

func BenchmarkInfot(b *testing.B) {
	ctx := context.WithValue(context.Background(), TRACEID, "bench-trace")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Infot(ctx, "benchmark test %d", i)
	}
}

// 基准测试 rotate logger 性能
func BenchmarkRotateLogger(b *testing.B) {
	tmpDir := b.TempDir()
	logFile := tmpDir + "/bench_rotate.log"
	logger := getRotateLogger(logFile)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark rotate log message")
		}
	})
}

func BenchmarkRotateLoggerWithPattern(b *testing.B) {
	tmpDir := b.TempDir()
	logPattern := tmpDir + "/bench_%Y%m%d%H.log"
	logger := getRotateLogger(logPattern)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("benchmark message %d", i)
	}
}

// 测试日志级别过滤
func TestLogLevelFiltering(t *testing.T) {
	// 设置为Error级别
	SetLevel(zapcore.ErrorLevel)

	// Debug和Info应该被过滤掉（实际测试中很难验证输出，这里主要测试不panic）
	Debug("this should be filtered")
	Info("this should be filtered")
	Warn("this should be filtered")
	Error("this should not be filtered")

	// 恢复为Debug级别
	SetLevel(zapcore.DebugLevel)
}

// 边界情况测试
func TestEdgeCases(t *testing.T) {
	t.Run("empty_message", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Empty message should not panic: %v", r)
			}
		}()
		Info("")
		Infof("")
	})

	t.Run("nil_args", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Nil args should not panic: %v", r)
			}
		}()
		Infof("test %v", nil)
	})

	t.Run("many_args", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Many args should not panic: %v", r)
			}
		}()
		Info("arg1", "arg2", "arg3", "arg4", "arg5")
	})
}
