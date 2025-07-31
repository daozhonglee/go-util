package xss

import (
	"strings"
	"testing"
)

func TestClean(t *testing.T) {
	// 测试清理XSS内容
	dirty := "<script>alert('xss')</script><p>正常内容</p>"
	clean := Clean(dirty)

	// 应该移除script标签但保留p标签
	if strings.Contains(clean, "<script>") {
		t.Error("Expected script tag to be removed")
	}

	if !strings.Contains(clean, "正常内容") {
		t.Error("Expected normal content to be preserved")
	}
}

func TestCleanEmpty(t *testing.T) {
	// 测试空内容
	result := Clean("")
	if result != "" {
		t.Error("Expected empty string for empty input")
	}
}
