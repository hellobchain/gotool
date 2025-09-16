package gtable

import (
	"bytes"
	"fmt"
	"strings"
)

type Table struct {
	headers []string
	rows    [][]string
}

func New() *Table {
	return &Table{}
}

func (t *Table) SetHeader(h ...string) *Table {
	t.headers = h
	return t
}

func (t *Table) AddRow(r ...interface{}) *Table {
	row := make([]string, len(r))
	for i, v := range r {
		row[i] = fmt.Sprint(v)
	}
	t.rows = append(t.rows, row)
	return t
}

func (t *Table) String() string {
	if len(t.headers) == 0 && len(t.rows) == 0 {
		return ""
	}
	colNum := len(t.headers)
	if colNum == 0 {
		colNum = len(t.rows[0])
	}
	widths := make([]int, colNum)
	for i := 0; i < colNum; i++ {
		if i < len(t.headers) {
			widths[i] = len(t.headers[i])
		}
		for _, row := range t.rows {
			if i < len(row) && len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
	}
	var buf bytes.Buffer
	// header
	if len(t.headers) > 0 {
		for i, h := range t.headers {
			buf.WriteString("| " + pad(h, widths[i]) + " ")
		}
		buf.WriteString("|\n")
		for i := 0; i < colNum; i++ {
			buf.WriteString("|-" + strings.Repeat("-", widths[i]) + "-")
		}
		buf.WriteString("|\n")
	}
	// rows
	for _, r := range t.rows {
		for i := 0; i < colNum; i++ {
			v := ""
			if i < len(r) {
				v = r[i]
			}
			buf.WriteString("| " + pad(v, widths[i]) + " ")
		}
		buf.WriteString("|\n")
	}
	return buf.String()
}

func (t *Table) Print() {
	fmt.Print(t.String())
}

func pad(s string, w int) string {
	if len(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-len(s))
}
