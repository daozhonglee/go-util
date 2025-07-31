package times

import (
	"testing"
	"time"
)

func TestGetCurrentUnix(t *testing.T) {
	unix := GetCurrentUnix()
	if unix <= 0 {
		t.Error("Expected positive unix timestamp")
	}
}

func TestGetCurrentMilliUnix(t *testing.T) {
	milli := GetCurrentMilliUnix()
	if milli <= 0 {
		t.Error("Expected positive millisecond timestamp")
	}
}

func TestDateStringFromTimestampMilli(t *testing.T) {
	// 2023-01-01 00:00:00 UTC 的毫秒时间戳
	timestamp := int64(1672531200000)
	result := DateStringFromTimestampMilli(timestamp, 0)
	expected := "2023-01-01"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestTimestampMilliToTime(t *testing.T) {
	// 测试毫秒时间戳转换为time.Time
	timestamp := int64(1672531200000) // 2023-01-01 00:00:00 UTC
	result := TimestampMilliToTime(timestamp)

	expected := time.Unix(1672531200, 0)
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestDateInt64FromTimestamp(t *testing.T) {
	// 测试时间戳转换为小时开始时间
	timestamp := int64(1672534800000) // 2023-01-01 01:00:00 UTC
	result := DateInt64FromTimestamp(timestamp)

	// 应该返回01:00:00的时间戳
	expected := int64(1672534800000)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
