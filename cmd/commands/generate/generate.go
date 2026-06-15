package generate

import (
	"fmt"
	"os"
	path "path/filepath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/zelviner/cgear/cmd/commands"
	"github.com/zelviner/cgear/cmd/commands/version"
	"github.com/zelviner/cgear/generate"
	"github.com/zelviner/cgear/logger"
	"github.com/zelviner/cgear/utils"
)

var CmdGenerate = &commands.Command{
	UsageLine: "generate [command]",
	Short:     "Generate project files or documentation",
	Long: `Creates an include directory or generates source.docx from the current project.

     $ cgear generate include
     $ cgear generate docs
	`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    GenerateCode,
}

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdGenerate)
}

func GenerateCode(cmd *commands.Command, args []string) int {
	currPath := utils.GetCgearWorkPath()
	if len(args) < 1 {
		logger.Log.Fatal("Command is missing")
	}

	gcmd := args[0]
	switch gcmd {
	case "include":
		include(cmd, args, currPath)

	case "docs":
		docs(cmd, args, currPath)

	default:
		logger.Log.Fatal("Command is missing")
	}

	logger.Log.Successf("%s successfully generated!", cases.Title(language.English).String(gcmd))

	return 0
}

func include(cmd *commands.Command, args []string, currPath string) {
	if !utils.IsCgearProject(currPath) {
		logger.Log.Fatal("not cgear project")
	}

	output := cmd.Out()

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(currPath, "include")+string(path.Separator), "\x1b[0m")
	os.MkdirAll(path.Join(currPath, "include"), 0755)

}

func docs(cmd *commands.Command, args []string, currPath string) {
	if !utils.IsCgearProject(currPath) {
		logger.Log.Fatal("not cgear project")
	}

	generate.SrcToDocx("source.docx", currPath)
}
