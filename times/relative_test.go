package times

import (
	"strings"
	"testing"
	"time"
)

func TestRelative(t *testing.T) {
	now := time.Now().Unix()

	// 测试刚刚
	result := Relative(now - 30)
	if result != "刚刚" {
		t.Errorf("Expected '刚刚', got %s", result)
	}

	// 测试分钟前
	result = Relative(now - 3*60)
	if !strings.Contains(result, "分钟前") {
		t.Errorf("Expected contains '分钟前', got %s", result)
	}

	// 测试小时前
	result = Relative(now - 2*3600)
	if !strings.Contains(result, "小时前") {
		t.Errorf("Expected contains '小时前', got %s", result)
	}
}
