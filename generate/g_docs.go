package generate

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"zel/logger"
	"zel/utils"
)

type FileInfo struct {
	FileName    string
	FileContent string
}

var FileTemplate = "```C++\n{{.FileContent}}\n```"

func GenerateMd(filename string, currPath string) {
	re := regexp.MustCompile(`\.(h|cpp|hpp)$`)

	var content = ""

	filepath.Walk(filepath.Join(currPath, "src"), func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if re.MatchString(info.Name()) {
			fc := trimFile(path)

			content += fc
		}
		return nil
	})

	utils.WriteToFile(filename, content)

}

func trimFile(filename string) string {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		logger.Log.Fatalf("无法打开文件:", err)
		return ""
	}
	defer file.Close()

	var connect string = ""

	// 逐行读取文件内容并去除空行
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			connect += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Log.Fatalf("扫描文件时出错:", err)
		return ""
	}

	return connect
}
