package generate

import (
	"os"
	"path/filepath"
	"regexp"
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
			fc := utils.FileTrim(path)

			content += fc
		}
		return nil
	})

	utils.WriteToFile(filename, content)

}
