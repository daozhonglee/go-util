package random

import "testing"

func TestInt(t *testing.T) {
	// 测试正常范围
	result := Int(1, 10)
	if result < 1 || result >= 10 {
		t.Errorf("Expected result in range [1,10), got %d", result)
	}
}

func TestIntInvalidRange(t *testing.T) {
	// 测试无效范围
	result := Int(10, 5)
	if result != 0 {
		t.Errorf("Expected 0 for invalid range, got %d", result)
	}
}

func TestIntEqualRange(t *testing.T) {
	// 测试相等范围
	result := Int(5, 5)
	if result != 0 {
		t.Errorf("Expected 0 for equal range, got %d", result)
	}
}
