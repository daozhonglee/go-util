package times

import (
	"testing"
	"time"
)

func TestTimeFormats(t *testing.T) {
	// 测试时间格式常量的正确性
	now := time.Now()

	// 测试DateTimeFormat
	dateTimeStr := now.Format(DateTimeFormat)
	parsedDateTime, err := time.Parse(DateTimeFormat, dateTimeStr)
	if err != nil {
		t.Errorf("DateTimeFormat parsing failed: %v", err)
	}

	// 验证格式是否正确（忽略毫秒和纳秒差异）
	if parsedDateTime.Year() != now.Year() ||
		parsedDateTime.Month() != now.Month() ||
		parsedDateTime.Day() != now.Day() ||
		parsedDateTime.Hour() != now.Hour() ||
		parsedDateTime.Minute() != now.Minute() ||
		parsedDateTime.Second() != now.Second() {
		t.Error("DateTimeFormat format validation failed")
	}

	// 测试DateFormat
	dateStr := now.Format(DateFormat)
	parsedDate, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		t.Errorf("DateFormat parsing failed: %v", err)
	}

	if parsedDate.Year() != now.Year() ||
		parsedDate.Month() != now.Month() ||
		parsedDate.Day() != now.Day() {
		t.Error("DateFormat format validation failed")
	}

	// 测试DateTimeMilliFormat
	dateTimeMilliStr := now.Format(DateTimeMilliFormat)
	parsedDateTimeMilli, err := time.Parse(DateTimeMilliFormat, dateTimeMilliStr)
	if err != nil {
		t.Errorf("DateTimeMilliFormat parsing failed: %v", err)
	}

	if parsedDateTimeMilli.Year() != now.Year() ||
		parsedDateTimeMilli.Month() != now.Month() ||
		parsedDateTimeMilli.Day() != now.Day() {
		t.Error("DateTimeMilliFormat format validation failed")
	}
}

func TestFormatValues(t *testing.T) {
	// 验证格式字符串的具体值
	if DateTimeFormat != "2006-01-02 15:04:05" {
		t.Errorf("Expected DateTimeFormat to be '2006-01-02 15:04:05', got '%s'", DateTimeFormat)
	}

	if DateTimeMilliFormat != "2006-01-02 15:04:05.000" {
		t.Errorf("Expected DateTimeMilliFormat to be '2006-01-02 15:04:05.000', got '%s'", DateTimeMilliFormat)
	}

	if DateFormat != "2006-01-02" {
		t.Errorf("Expected DateFormat to be '2006-01-02', got '%s'", DateFormat)
	}
}
