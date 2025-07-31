// Package random 提供随机数生成功能
package random

import (
	"math/rand"
	"time"
)

// 全局随机数生成器，程序启动时初始化一次
var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Int 生成区间[min,max)的随机整数
// 如果 min == max，返回该值
// 如果 min > max，返回0（无效范围）
func Int(min, max int) int {
	if min > max {
		return 0 // 无效范围
	}
	if min == max {
		return min // 相等时返回该值
	}
	return min + globalRand.Intn(max-min)
}
