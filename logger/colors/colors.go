package colors

import (
	"fmt"
	"io"
)

// 输出模式
type outputMode int

// DiscardNonColorEscSeq支持分割的颜色转义序列。
// 但不输出非彩色转义序列。
// 如果你想输出一个非颜色的转义序列，比如ncurses，请使用OutputNonColorEscSeq。但是，它不支持分割的颜色转义序列。
const (
	_ outputMode = iota
	DiscardNonColorEscSeq
	OutputNonColorEscSeq
)

// NewColorWriter使用io.Writer w作为其初始内容创建并重命名一个新的ansiColorWriter。
// 在Windows的控制台中，通过转义序列改变文本的前景色和背景色。
// 在其他系统的控制台中，写入w all text。
func NewColorWriter(w io.Writer) io.Writer {
	return NewModeColorWriter(w, DiscardNonColorEscSeq)
}

// NewModeColorWriter通过指定outputMode来创建并调用一个新的ansiColorWriter。
func NewModeColorWriter(w io.Writer, mode outputMode) io.Writer {
	if _, ok := w.(*colorWriter); !ok {
		return &colorWriter{
			w:    w,
			mode: mode,
		}
	}
	return w
}

// Bold 返回一个粗体字符串
func Bold(message string) string {
	return fmt.Sprintf("\x1b[1m%s\x1b[0m", message)
}

// Black 返回一个黑色字符串
func Black(message string) string {
	return fmt.Sprintf("\x1b[30m%s\x1b[0m", message)
}

// White 返回一个白色字符串
func White(message string) string {
	return fmt.Sprintf("\x1b[37m%s\x1b[0m", message)
}

// Cyan 返回一个青色字符串
func Cyan(message string) string {
	return fmt.Sprintf("\x1b[36m%s\x1b[0m", message)
}

// Blue 返回一个蓝色字符串
func Blue(message string) string {
	return fmt.Sprintf("\x1b[34m%s\x1b[0m", message)
}

// Red 返回一个红色字符串
func Red(message string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", message)
}

// Green 返回一个绿色字符串
func Green(message string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", message)
}

// Yellow 返回一个黄色字符串
func Yellow(message string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", message)
}

// Gray 返回一个灰色字符串
func Gray(message string) string {
	return fmt.Sprintf("\x1b[37m%s\x1b[0m", message)
}

// Magenta 返回一个洋红色字符串
func Magenta(message string) string {
	return fmt.Sprintf("\x1b[35m%s\x1b[0m", message)
}

// BlackBold 返回一个黑色粗体字符串
func BlackBold(message string) string {
	return fmt.Sprintf("\x1b[30m%s\x1b[0m", Bold(message))
}

// WhiteBold 返回一个白色粗体字符串
func WhiteBold(message string) string {
	return fmt.Sprintf("\x1b[37m%s\x1b[0m", Bold(message))
}

// CyanBold 返回一个青色粗体字符串
func CyanBold(message string) string {
	return fmt.Sprintf("\x1b[36m%s\x1b[0m", Bold(message))
}

// BlueBold 返回一个蓝色粗体字符串
func BlueBold(message string) string {
	return fmt.Sprintf("\x1b[34m%s\x1b[0m", Bold(message))
}

// RedBold 返回一个红色粗体字符串
func RedBold(message string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", Bold(message))
}

// GreenBold 返回一个绿色粗体字符串
func GreenBold(message string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", Bold(message))
}

// YellowBold 返回一个黄色粗体字符串
func YellowBold(message string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", Bold(message))
}

// GrayBold 返回一个灰色粗体字符串
func GrayBold(message string) string {
	return fmt.Sprintf("\x1b[37m%s\x1b[0m", Bold(message))
}

// MagentaBold 返回一个紫色粗体字符串
func MagentaBold(message string) string {
	return fmt.Sprintf("\x1b[35m%s\x1b[0m", Bold(message))
}
