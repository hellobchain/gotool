package gexcel

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ================= CSV =================
func WriteCSV(path string, headers []string, rows [][]interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	if err = w.Write(headers); err != nil {
		return err
	}
	for _, row := range rows {
		record := make([]string, len(row))
		for i, v := range row {
			record[i] = fmt.Sprintf("%v", v)
		}
		if err = w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// ================= XLSX =================
// 只写入单个工作表，零依赖
func WriteXLSX(path string, headers []string, rows [][]interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	// 1. [Content_Types].xml
	if err := addZipFile(zw, "[Content_Types].xml", contentTypes); err != nil {
		return err
	}
	// 2. _rels/.rels
	if err := addZipFile(zw, "_rels/.rels", rels); err != nil {
		return err
	}
	// 3. xl/workbook.xml
	if err := addZipFile(zw, "xl/workbook.xml", workbook); err != nil {
		return err
	}
	// 4. xl/_rels/workbook.xml.rels
	if err := addZipFile(zw, "xl/_rels/workbook.xml.rels", workbookRels); err != nil {
		return err
	}
	// 5. xl/worksheets/sheet1.xml
	sheetXML := genSheetXML(headers, rows)
	if err := addZipFile(zw, "xl/worksheets/sheet1.xml", sheetXML); err != nil {
		return err
	}
	// 6. xl/styles.xml（最小样式）
	if err := addZipFile(zw, "xl/styles.xml", styles); err != nil {
		return err
	}
	return zw.Close()
}

func addZipFile(zw *zip.Writer, name, content string) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(content))
	return err
}

// 生成工作表 XML
func genSheetXML(headers []string, rows [][]interface{}) string {
	var buf bytes.Buffer
	buf.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">`)
	buf.WriteString(`<sheetData>`)
	// headers
	writeRow(&buf, 1, headers)
	// data
	for i, row := range rows {
		cells := make([]string, len(row))
		for j, v := range row {
			cells[j] = fmt.Sprintf("%v", v)
		}
		writeRow(&buf, i+2, cells)
	}
	buf.WriteString(`</sheetData></worksheet>`)
	return buf.String()
}

func writeRow(buf *bytes.Buffer, rowIndex int, cells []string) {
	buf.WriteString(fmt.Sprintf(`<row r="%d">`, rowIndex))
	for i, v := range cells {
		// 简单判断数字
		cell := fmt.Sprintf(`<c r="%s%s" t="str"><v>%s</v></c>`,
			colName(i), strconv.Itoa(rowIndex), escape(v))
		if _, err := strconv.ParseFloat(v, 64); err == nil {
			cell = fmt.Sprintf(`<c r="%s%s"><v>%s</v></c>`, colName(i), strconv.Itoa(rowIndex), v)
		}
		buf.WriteString(cell)
	}
	buf.WriteString(`</row>`)
}

func colName(i int) string {
	name := ""
	for i >= 0 {
		name = string(rune('A'+i%26)) + name
		i = i/26 - 1
	}
	return name
}

func escape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// ========== 最小 OOXML 模板 ==========
const contentTypes = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
  <Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
  <Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/>
</Types>`

const rels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`

const workbookRels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
</Relationships>`

const workbook = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
  <sheets>
    <sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/>
  </sheets>
</workbook>`

const styles = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
  <fonts count="1"><font><sz val="11"/><name val="Calibri"/></font></fonts>
  <fills count="1"><fill><patternFill patternType="none"/></fill></fills>
  <borders count="1"><border><left/><right/><top/><bottom/></border></borders>
  <cellXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0"/></cellXfs>
</styleSheet>`
