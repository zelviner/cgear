package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestFilePath(t *testing.T) {
	// 获取当前文件夹路径
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 遍历当前文件夹
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 检查错误
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		// 忽略当前文件夹
		if path != dir {
			// 判断是否为文件夹
			if info.IsDir() {
				fmt.Println("文件夹:", info.Name())
			} else {
				fmt.Println("文件:", info.Name())
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func TestPathSplit(t *testing.T) {

	str := "C:/a/b/c/d/e; /sdf/sdf/aaa; /cc/gg/ee"

	strList := filepath.SplitList(str)

	for _, str := range strList {
		fmt.Println(str)
	}

}
