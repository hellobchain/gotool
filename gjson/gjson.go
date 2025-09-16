package gjson

import (
	"encoding/json"
	"os"
)

// ToStringStruct 任意结构转 JSON 字符串
func ToStringStruct(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// ToStringMap 转 map 再 JSON
func ToStringMap(m map[string]interface{}) string { return ToStringStruct(m) }

// ParseString 解析 JSON 字符串到 map
func ParseString(str string) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(str), &m)
	return m, err
}

// ParseFile 解析 JSON 文件到 map
func ParseFile(path string) (map[string]interface{}, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseString(string(b))
}
