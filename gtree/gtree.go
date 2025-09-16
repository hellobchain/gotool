package gtree

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Option 配置项
type Option struct {
	// 仅目录
	DirOnly bool
	// 显示文件大小
	ShowSize bool
	// 最大深度（<0 表示不限制）
	MaxDepth int
	// 自定义过滤器，返回 true 则跳过
	Skip func(path string, d fs.DirEntry) bool
}

// DefaultOption 默认配置
var DefaultOption = Option{
	DirOnly:  false,
	ShowSize: false,
	MaxDepth: -1,
}

// Print 打印当前工作目录的树
func Print() error {
	return PrintDir(".", DefaultOption)
}

// PrintDir 打印指定目录的树
func PrintDir(root string, opt ...Option) error {
	option := DefaultOption
	if len(opt) > 0 {
		option = opt[0]
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	fmt.Println(abs)
	return walk(abs, "", 0, option)
}

// walk 递归遍历
func walk(root, prefix string, depth int, opt Option) error {
	if opt.MaxDepth >= 0 && depth >= opt.MaxDepth {
		return nil
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	// 过滤
	var list []fs.DirEntry
	for _, e := range entries {
		if opt.Skip != nil && opt.Skip(filepath.Join(root, e.Name()), e) {
			continue
		}
		if opt.DirOnly && !e.IsDir() {
			continue
		}
		list = append(list, e)
	}
	for i, e := range list {
		isLast := i == len(list)-1
		name := e.Name()
		if opt.ShowSize && !e.IsDir() {
			info, _ := e.Info()
			name += fmt.Sprintf(" (%d)", info.Size())
		}
		fmt.Print(prefix)
		if isLast {
			fmt.Print("└── ")
		} else {
			fmt.Print("├── ")
		}
		fmt.Println(name)
		// 下一层前缀
		nextPrefix := prefix
		if isLast {
			nextPrefix += "    "
		} else {
			nextPrefix += "│   "
		}
		if e.IsDir() {
			if err := walk(filepath.Join(root, e.Name()), nextPrefix, depth+1, opt); err != nil {
				// 允许继续打印兄弟目录
			}
		}
	}
	return nil
}

// String 返回树字符串（不直接打印）
func String(root string, opt ...Option) (string, error) {
	var sb strings.Builder
	option := DefaultOption
	if len(opt) > 0 {
		option = opt[0]
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	sb.WriteString(abs + "\n")
	if err := walkString(abs, "", 0, option, &sb); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func walkString(root, prefix string, depth int, opt Option, sb *strings.Builder) error {
	if opt.MaxDepth >= 0 && depth >= opt.MaxDepth {
		return nil
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	var list []fs.DirEntry
	for _, e := range entries {
		if opt.Skip != nil && opt.Skip(filepath.Join(root, e.Name()), e) {
			continue
		}
		if opt.DirOnly && !e.IsDir() {
			continue
		}
		list = append(list, e)
	}
	for i, e := range list {
		isLast := i == len(list)-1
		name := e.Name()
		if opt.ShowSize && !e.IsDir() {
			info, _ := e.Info()
			name += fmt.Sprintf(" (%d)", info.Size())
		}
		sb.WriteString(prefix)
		if isLast {
			sb.WriteString("└── ")
		} else {
			sb.WriteString("├── ")
		}
		sb.WriteString(name + "\n")
		nextPrefix := prefix
		if isLast {
			nextPrefix += "    "
		} else {
			nextPrefix += "│   "
		}
		if e.IsDir() {
			_ = walkString(filepath.Join(root, e.Name()), nextPrefix, depth+1, opt, sb)
		}
	}
	return nil
}
