// Package json 提供便捷的JSON序列化和格式化功能
package json

import (
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
)

// Marshal 将对象序列化为JSON字符串，出错时返回空字符串
func MarshalFailSafe(value interface{}) string {
	result, err := sonic.Marshal(value)
	if err != nil {
		return ""
	}
	return string(result)
}

// Clean 清理JSON字符串，移除多余的字符和格式符号
func Clean(s string) string {
	s = strings.ReplaceAll(s, "Json", "")
	s = strings.ReplaceAll(s, "json", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "```", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\\t", "")
	s = strings.ReplaceAll(s, "\\n", "")
	s = strings.ReplaceAll(s, "\\", "")
	s = strings.ReplaceAll(s, "\"\"\"", "")
	s = strings.ReplaceAll(s, "æ", "")
	return s
}

// Pretty 格式化打印JSON数据
func PrettyFailSafe(data interface{}) {
	jsonData, err := sonic.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(jsonData))
}

// Unmarshal JSON反序列化
func Unmarshal(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}
