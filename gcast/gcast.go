package gcast

import (
	"encoding/hex"
	"strconv"
)

func ToString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return strconv.FormatInt(ToInt64(v), 10)
	case float32, float64:
		return strconv.FormatFloat(ToFloat64(v), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}

func ToHex(v interface{}) string {
	switch v := v.(type) {
	case string:
		return hex.EncodeToString([]byte(v))
	case []byte:
		return hex.EncodeToString(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return strconv.FormatInt(ToInt64(v), 16)
	case float32, float64:
		return strconv.FormatFloat(ToFloat64(v), 'f', -1, 64)
	default:
		return ""
	}
}

func ToInt(v interface{}) int { return int(ToInt64(v)) }
func ToInt64(v interface{}) int64 {
	switch v := v.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	default:
		return 0
	}
}
func ToFloat64(v interface{}) float64 {
	switch v := v.(type) {
	case float64:
		return v
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0
	}
}
func ToBool(v interface{}) bool {
	switch v := v.(type) {
	case bool:
		return v
	case string:
		b, _ := strconv.ParseBool(v)
		return b
	default:
		return false
	}
}
