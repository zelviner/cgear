package generate

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"zel/logger"

	"github.com/ErmaiSoft/GoOpenXml/word"
)

type FileInfo struct {
	FileName    string
	FileContent string
}

var FileTemplate = "```C++\n{{.FileContent}}\n```"

func SrcToDocx(filename string, currPath string) {

	paragraph := GetParagraph(currPath)

	WriteToDocx(filename, paragraph)

}

func GetParagraph(currPath string) (paragraph []word.Paragraph) {
	re := regexp.MustCompile(`\.(h|cpp|hpp)$`)
	filepath.Walk(filepath.Join(currPath, "src"), func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if re.MatchString(info.Name()) {
			// 打开文件
			file, err := os.Open(path)
			if err != nil {
				logger.Log.Fatalf("无法打开文件: %s", err)
			}
			defer file.Close()

			font := word.Font{Family: "Consolas", Size: 10, Bold: false, Color: "000000"} //字体
			lineSeting := word.Line{Rule: word.LineRuleAuto, Height: 1}                   //行高、行间距、首行缩进

			// 逐行读取文件内容
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if line != "" {
					paragraph = append(paragraph, word.Paragraph{
						F: font,
						L: lineSeting,
						T: []word.Text{
							{T: line, F: &font},
						},
					})
				}
			}

			if err := scanner.Err(); err != nil {
				logger.Log.Fatalf("扫描文件时出错: %s", err)
			}
		}
		return nil
	})

	return
}

// 写入内容到 .docx 文件
func WriteToDocx(filename string, paragraph []word.Paragraph) {
	pos := strings.IndexByte(filename, '.')
	suffix := filename[pos:]
	if suffix != ".docx" {
		logger.Log.Fatalf("File name suffix error, need '.docx', get '%s'", suffix)
	}

	docx := word.CreateDocx()
	docx.AddParagraph(paragraph)

	err := docx.WriteToFile(filename)
	if err != nil {
		logger.Log.Error(err.Error())
	}

}
