package ghttp

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var defaultClient = &http.Client{Timeout: 10 * time.Second}

// Get 快速 GET 请求，返回 body 字符串
func Get(url string) (string, error) {
	resp, err := defaultClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

// PostJSON 发送 JSON 并返回 body 字符串
func PostJSON(url string, obj interface{}) (string, error) {
	raw, _ := json.Marshal(obj)
	resp, err := defaultClient.Post(url, "application/json", bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

// post form data
func PostForm(url string, data url.Values) (string, error) {
	resp, err := defaultClient.PostForm(url, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

// UploadFile
type UploadFile struct {
	FieldName string // 表单字段名
	FileName  string // 远程文件名（可空，默认用本地名）
	Reader    io.Reader
}

func PostMultipart(url string, files []UploadFile, fields map[string]string) (string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// 写入文件
	for _, f := range files {
		// 如果 Reader 是 *os.File，则尽量取真实文件名
		realName := f.FileName
		if realName == "" {
			if rf, ok := f.Reader.(*os.File); ok {
				realName = filepath.Base(rf.Name())
			} else {
				realName = "unknown"
			}
		}
		part, err := w.CreateFormFile(f.FieldName, realName)
		if err != nil {
			return "", err
		}
		if _, err = io.Copy(part, f.Reader); err != nil {
			return "", err
		}
	}

	// 写入普通字段
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return "", err
		}
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	// 发送请求
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := defaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func PostFile(url, fieldName string, file io.Reader, extraFields map[string]string) (string, error) {
	return PostMultipart(url, []UploadFile{{FieldName: fieldName, Reader: file}}, extraFields)
}
