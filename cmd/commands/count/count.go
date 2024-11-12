package count

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/utils"
)

var CmdCount = &commands.Command{
	UsageLine: "count",
	Short:     "Counting source file lines",
	Long: `
Counting source file (.h, .hpp, .cpp) lines.`,
	Run: Count,
}

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdCount)
}

func Count(cmd *commands.Command, args []string) int {

	projectPath := utils.GetZelWorkPath()

	var lines int = 0
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".h" || filepath.Ext(path) == ".hpp" || filepath.Ext(path) == ".cpp" {
			tempLine, err := countLines(path)
			if err != nil {
				return err
			}
			lines += tempLine
		}

		return nil
	}

	filepath.Walk(projectPath, walkFunc)

	logger.Log.Successf("lines: %d", lines)

	return 0
}

// countLines 统计文件的行数
func countLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}
