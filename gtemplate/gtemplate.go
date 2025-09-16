package gtemplate

import (
	"bytes"
	"fmt"
	"strings"
)

// Render 仅支持 {{.Key}} 简单替换
func Render(tpl string, data map[string]interface{}) string {
	var buf bytes.Buffer
	i := 0
	for i < len(tpl) {
		if i < len(tpl)-3 && tpl[i:i+3] == "{{." {
			j := strings.Index(tpl[i+3:], "}}")
			if j >= 0 {
				key := tpl[i+3 : i+3+j]
				if v, ok := data[key]; ok {
					buf.WriteString(strings.ReplaceAll(fmt.Sprint(v), "{{", "\\{{"))
				}
				i += 3 + j + 2
				continue
			}
		}
		buf.WriteByte(tpl[i])
		i++
	}
	return buf.String()
}
