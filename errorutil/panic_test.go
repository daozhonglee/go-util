package errorutil

import (
	"errors"
	"testing"
)

func TestPanicIf(t *testing.T) {
	// 测试nil错误不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Error("Expected no panic for nil error")
		}
	}()

	PanicIf(nil)
}

func TestPanicIfWithError(t *testing.T) {
	// 测试非nil错误会panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for non-nil error")
		}
	}()

	PanicIf(errors.New("test error"))
}

func TestPanicWithStack(t *testing.T) {
	// 测试带堆栈的panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with stack")
		}
	}()

	PanicIfWithStack(errors.New("test error"))
}
