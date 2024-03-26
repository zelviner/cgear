package tests

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	// 待匹配的文件列表
	files := []string{"file.h", "file.cpp", "file.hpp", "file.txt", "file"}

	// 编译正则表达式
	re := regexp.MustCompile(`\.(h|cpp|hpp)$`)

	// 遍历文件列表，匹配正则表达式
	for _, file := range files {
		if re.MatchString(file) {
			fmt.Printf("%s 匹配\n", file)
		} else {
			fmt.Printf("%s 不匹配\n", file)
		}
	}
}
