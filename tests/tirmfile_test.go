package tests

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestTirmFile(t *testing.T) {
	// 打开文件
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	// 创建一个新文件用于写入非空行内容
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("无法创建输出文件:", err)
		return
	}
	defer outputFile.Close()

	// 逐行读取文件内容并去除空行
	scanner := bufio.NewScanner(file)
	var connect string = ""
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			connect += line + "\n"
		}
	}

	outputFile.WriteString(connect)

	if err := scanner.Err(); err != nil {
		fmt.Println("扫描文件时出错:", err)
		return
	}

	fmt.Println("空行已经被去除并已写入到 output.txt 文件中.")
}
