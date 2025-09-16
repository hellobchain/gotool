package gstr

import (
	"strings"
	"unicode"
)

func IsBlank(s string) bool    { return strings.TrimSpace(s) == "" }
func IsEmpty(s string) bool    { return s == "" }
func IsNotBlank(s string) bool { return !IsBlank(s) }
func IsNotEmpty(s string) bool { return !IsEmpty(s) }

// SubBetween 提取 start 与 end 之间的文本
func SubBetween(s, start, end string) string {
	a := strings.Index(s, start)
	if a == -1 {
		return ""
	}
	b := strings.Index(s[a+len(start):], end)
	if b == -1 {
		return ""
	}
	return s[a+len(start) : a+len(start)+b]
}

// Reverse 字符串反转
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// CamelCase 下划线转小驼峰
func CamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// SnakeCase 驼峰转下划线
func SnakeCase(s string) string {
	var res []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			res = append(res, '_')
		}
		res = append(res, unicode.ToLower(r))
	}
	return string(res)
}
