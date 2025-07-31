// Package id 提供身份证相关的验证工具
package id

import (
	"fmt"
	"regexp"
	"time"
)

// ValidateChinese 验证中国18位身份证号码格式，支持末位X
func ValidateChinese(id string) bool {
	// 18位身份证号码的正则表达式
	pattern := `^\d{17}[\dxX]$`
	matched, err := regexp.MatchString(pattern, id)
	if err != nil {
		fmt.Println("Error matching regex:", err)
		return false
	}

	return matched
}

// CheckAge16To18 判断年龄是否在16-18岁范围（包含16岁，不包含18岁）
func CheckAge16To18(id string) bool {
	if !ValidateChinese(id) {
		return false
	}
	// 假设身份证号码中的出生日期格式为 "20060102"（年月日）
	dobStr := id[6:14]
	dob, err := time.Parse("20060102", dobStr)
	if err != nil {
		fmt.Println("Error parsing date of birth:", err)
		return false
	}

	// 计算年龄
	age := time.Now().Year() - dob.Year()

	return age >= 16 && age < 18
}
