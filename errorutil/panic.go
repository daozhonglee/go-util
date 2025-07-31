// Package errorutil 提供错误处理和panic恢复功能
package errorutil

import (
	"fmt"
	"runtime/debug"
)

// PanicIf 如果错误不为nil则panic
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// PanicWithStack 如果错误不为nil则panic，并包含堆栈信息
func PanicIfWithStack(err error) {
	if err != nil {
		panic(fmt.Sprintf("error: %v\nstack: %v", err, string(debug.Stack())))
	}
}

// Recover 通用panic恢复函数
func Recover() {
	if err := recover(); err != nil {
		fmt.Printf("recover panic, err = %v\n", err)
	}
}
