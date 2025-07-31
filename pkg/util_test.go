package util

import (
	"testing"
)

func TestMainEntryAlias(t *testing.T) {
	// 测试ID验证别名
	result := ValidateChineseID("11010119900307001X")
	if !result {
		t.Error("ValidateChineseID alias should work")
	}

	// 测试MD5别名
	hash := MD5Hash([]byte("test"))
	expected := "098f6bcd4621d373cade4e832627b4f6"
	if hash != expected {
		t.Errorf("MD5Hash alias should work, expected %s, got %s", expected, hash)
	}

	// 测试JSON别名
	jsonStr := JSONMarshal(map[string]string{"key": "value"})
	if jsonStr == "" {
		t.Error("JSONMarshal alias should work")
	}

	// 测试随机数别名
	num := RandomInt(1, 10)
	if num < 1 || num >= 10 {
		t.Errorf("RandomInt alias should work, got %d", num)
	}

	// 测试集合别名
	set := NewSet(1, 2, 3)
	if set.Len() != 3 {
		t.Error("NewSet alias should work")
	}
}
