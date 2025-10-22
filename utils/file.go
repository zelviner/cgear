package utils

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zelviner/cgear/logger"
)

// 检查文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

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
	desFile, err := os.OpenFile(des, os.O_RDWR|os.O_CREATE, perm) //复制源文件的所有权限
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
	var (
		file *os.File
		err  error
	)
	if !IsExist(filename) {
		file, err = os.Create(filename)
	} else {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	}
	MustCheck(err)
	defer CloseFile(file)

	_, err = file.WriteString(content)
	MustCheck(err)
}

func ReadFile(filename string) string {
	file, err := os.Open(filename)
	MustCheck(err)
	defer CloseFile(file)

	bytes, err := io.ReadAll(file)
	MustCheck(err)

	return string(bytes)
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

func ReplaceFileContent(filename string, old string, new string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer CloseFile(file)

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	content := strings.Replace(string(bytes), old, new, -1)
	WriteToFile(filename, content)

	return err
}

// ZipFile compresses the specified source directory and saves the result as a zip file at the destination path.
func ZipFile(src, dst string) error {
	// Create the zip file.
	zipFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	// Create a new zip archive.
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the source directory and add files to the zip archive.
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path of the file or directory to be used as the header name.
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Create a zip file header based on the file info.
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// Create a writer for the file header and copy the file content.
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk through source directory: %w", err)
	}

	return nil
}
