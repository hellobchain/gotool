package gcli

import (
	"fmt"
)

const (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

func Red(s string) string     { return red + s + reset }
func Green(s string) string   { return green + s + reset }
func Yellow(s string) string  { return yellow + s + reset }
func Blue(s string) string    { return blue + s + reset }
func Magenta(s string) string { return magenta + s + reset }
func Cyan(s string) string    { return cyan + s + reset }

// PrintSuccess 绿色对勾
func PrintSuccess(format string, a ...interface{}) {
	fmt.Println(Green("✔ " + fmt.Sprintf(format, a...)))
}

// PrintError 红色叉
func PrintError(format string, a ...interface{}) {
	fmt.Println(Red("✖ " + fmt.Sprintf(format, a...)))
}
