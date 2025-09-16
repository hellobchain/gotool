package gcolor

import (
	"fmt"
	"strconv"
)

const esc = "\033["

type Color string

const (
	Reset        Color = "0"
	Bold         Color = "1"
	RED          Color = "31"
	GREEN        Color = "32"
	YELLOW       Color = "33"
	BLUE         Color = "34"
	MAGENTA      Color = "35"
	CYAN         Color = "36"
	LightGray    Color = "37"
	DarkGray     Color = "90"
	LightRed     Color = "91"
	LightGreen   Color = "92"
	LightYellow  Color = "93"
	LightBlue    Color = "94"
	LightMagenta Color = "95"
	LightCyan    Color = "96"
)

func colorize(c Color, v ...interface{}) string {
	return esc + string(c) + "m" + fmt.Sprint(v...) + esc + string(Reset) + "m"
}

func Red(v ...interface{}) string     { return colorize(RED, v...) }
func Green(v ...interface{}) string   { return colorize(GREEN, v...) }
func Yellow(v ...interface{}) string  { return colorize(YELLOW, v...) }
func Blue(v ...interface{}) string    { return colorize(BLUE, v...) }
func Magenta(v ...interface{}) string { return colorize(MAGENTA, v...) }
func Cyan(v ...interface{}) string    { return colorize(CYAN, v...) }

// RGB 真彩色
func RGB(r, g, b int, v ...interface{}) string {
	return esc + "38;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b) + "m" +
		fmt.Sprint(v...) + esc + string(Reset) + "m"
}
