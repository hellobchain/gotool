package gfile

import (
	"os"
	"path/filepath"
)

// ReadString 一次性读取文本文件
func ReadString(path string) (string, error) {
	b, err := os.ReadFile(path)
	return string(b), err
}

// WriteString 一次性写入文本文件（覆盖）
func WriteString(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// 追加
func AppendString(path string, content string) error {
	// 在原先文件基础上追加内容
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

// Exists 文件或目录是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// IsDir 是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// MkdirIfNot 目录不存在则创建
func MkdirIfNot(path string) error {
	if !Exists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// Ext 获取文件扩展名（含点）
func Ext(path string) string { return filepath.Ext(path) }
