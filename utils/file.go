package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZEL-30/zel/logger"
)

// // 检查文件是否存在
// func FileIsExisted(filename string) bool {
// 	existed := true
// 	if _, err := os.Stat(filename); os.IsNotExist(err) {
// 		existed = false
// 	}
// 	return existed
// }

// 创建文件夹（如果文件夹不存在则创建）
func MakeDir(dir string) error {
	if !IsExist(dir) {
		if err := os.MkdirAll(dir, 0777); err != nil { //os.ModePerm
			fmt.Println("MakeDir failed:", err)
			return err
		}
	}
	return nil
}

// 复制文件
func CopyFile(src, des string) (written int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer CloseFile(srcFile)

	//获取源文件的权限
	fi, _ := srcFile.Stat()
	perm := fi.Mode()

	//desFile, err := os.Create(des)  //无法复制源文件的所有权限
	desFile, err := os.OpenFile(des, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm) //复制源文件的所有权限
	if err != nil {
		return 0, err
	}
	defer CloseFile(desFile)

	return io.Copy(desFile, srcFile)
}

// 复制文件夹
func CopyDir(srcPath, desPath string) error {
	//检查目录是否正确
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			return errors.New("源路径不是一个正确的目录！")
		}
	}

	MakeDir(desPath)

	if strings.TrimSpace(srcPath) == strings.TrimSpace(desPath) {
		return errors.New("源路径与目标路径不能相同！")
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		//复制目录是将源目录中的子目录复制到目标路径中，不包含源目录本身
		if path == srcPath {
			return nil
		}

		//生成新路径
		destNewPath := strings.Replace(path, srcPath, desPath, -1)

		if !f.IsDir() {
			CopyFile(path, destNewPath)
		} else {
			if !IsExist(destNewPath) {
				return MakeDir(destNewPath)
			}
		}

		return nil
	})

	return err
}

// 创建文件并向其中写入内容
func WriteToFile(filename string, content string) {
	f, err := os.Create(filename)
	MustCheck(err)
	defer CloseFile(f)

	_, err = f.WriteString(content)
	MustCheck(err)
}

// 尝试关闭传递的文件, 如果出错 panic
func CloseFile(f *os.File) {
	err := f.Close()
	MustCheck(err)
}

// 去除文件中的空行
func FileTrim(filename string) (content string) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		logger.Log.Fatalf("无法打开文件: %s", err)
		return
	}
	defer CloseFile(file)

	// 逐行读取文件内容并去除空行
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			content += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Log.Fatalf("扫描文件时出错: %s", err)
		return
	}

	return
}
