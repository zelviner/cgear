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
	return err == nil
}

// 创建文件夹（如果文件夹不存在则创建）
func MakeDir(dir string) error {
	if !IsExist(dir) {
		return os.MkdirAll(dir, 0755)
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

	fi, err := srcFile.Stat()
	if err != nil {
		return 0, err
	}
	perm := fi.Mode()

	desFile, err := os.OpenFile(des, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
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

	if strings.TrimSpace(srcPath) == strings.TrimSpace(desPath) {
		return errors.New("源路径与目标路径不能相同！")
	}
	if err := MakeDir(desPath); err != nil {
		return err
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//复制目录是将源目录中的子目录复制到目标路径中，不包含源目录本身
		if path == srcPath {
			return nil
		}

		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}
		destNewPath := filepath.Join(desPath, relPath)

		if !f.IsDir() {
			_, err := CopyFile(path, destNewPath)
			return err
		}
		if !IsExist(destNewPath) {
			return MakeDir(destNewPath)
		}

		return nil
	})

	return err
}

// 创建文件并向其中写入内容
func WriteToFile(filename string, content string) {
	MustCheck(os.WriteFile(filename, []byte(content), 0644))
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

	content := strings.ReplaceAll(string(bytes), old, new)
	return os.WriteFile(filename, []byte(content), 0644)
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

			_, copyErr := io.Copy(writer, file)
			closeErr := file.Close()
			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}
		}

		return nil
	})

	if err != nil {
		_ = zipWriter.Close()
		return fmt.Errorf("failed to walk through source directory: %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close zip file: %w", err)
	}

	return nil
}
