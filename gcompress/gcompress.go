package gcompress

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// GzipBytes gzip 压缩字节
func GzipBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GunzipBytes gzip 解压
func GunzipBytes(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(r)
}

// ZipFile 压缩单个文件或目录
func ZipFile(src, dst string) error {
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return filepath.Walk(src, func(path string, fi os.FileInfo, err error) error {
			if err != nil || fi.IsDir() {
				return err
			}
			rel, _ := filepath.Rel(src, path)
			fw, err := w.Create(rel)
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(fw, file)
			return err
		})
	}
	fw, err := w.Create(filepath.Base(src))
	if err != nil {
		return err
	}
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(fw, file)
	return err
}
