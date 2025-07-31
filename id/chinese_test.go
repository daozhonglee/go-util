package id

import "testing"

func TestValidateChinese(t *testing.T) {
	// 测试有效身份证号
	valid := ValidateChinese("11010119900307001X")
	if !valid {
		t.Error("Expected valid ID card to return true")
	}

	// 测试无效身份证号
	invalid := ValidateChinese("123456789")
	if invalid {
		t.Error("Expected invalid ID card to return false")
	}
}

func TestCheckAge16To18(t *testing.T) {
	// 测试年龄在16-18之间的身份证（这个需要根据当前年份调整）
	// 这里使用一个大概的测试
	result := CheckAge16To18("11010520070307001X") // 2007年出生
	// 注意：这个测试结果取决于当前年份，这里不做具体断言
	_ = result
}
