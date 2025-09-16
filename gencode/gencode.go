package gencode

import (
	"encoding/base64"
	"encoding/hex"
)

var base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Base62 编码整数
func Base62(n uint64) string {
	if n == 0 {
		return "0"
	}
	var s []byte
	for n > 0 {
		s = append(s, base62Alphabet[n%62])
		n /= 62
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return string(s)
}

// URLSafe base64 URL 安全
func URLSafe(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}

// Hex hex 编码
func Hex(s string) string { return hex.EncodeToString([]byte(s)) }
