package async

import (
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	future := Go(func() interface{} {
		return "test"
	})

	result := future.Get()
	if result != "test" {
		t.Errorf("Expected 'test', got %v", result)
	}
}

func TestTimeout(t *testing.T) {
	// 测试正常执行
	result := Timeout(func() interface{} {
		return "success"
	}, 100*time.Millisecond)

	if result != "success" {
		t.Errorf("Expected 'success', got %v", result)
	}

	// 测试超时
	result = Timeout(func() interface{} {
		time.Sleep(200 * time.Millisecond)
		return "timeout"
	}, 50*time.Millisecond)

	if result != nil {
		t.Error("Expected nil for timeout")
	}
}

func TestSafe(t *testing.T) {
	// 测试安全执行（不应该panic）
	Safe(func() {
		panic("test panic")
	})

	// 如果到这里说明没有panic，测试通过
	time.Sleep(10 * time.Millisecond) // 等待goroutine执行
}
