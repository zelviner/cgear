package new

import (
	"fmt"
	"os"
	path "path/filepath"
	"strings"

	"zel/cmd/commands"
	"zel/cmd/commands/version"
	"zel/logger"
	"zel/logger/colors"
	"zel/utils"
)

var qt utils.DocValue
var zelVersion utils.DocValue

var CmdNew = &commands.Command{
	UsageLine: "new [appname] [-qt=false]",
	Short:     "Create a C++ project, using cmake tool and opening by vscode",
	Long: `Creates a C++ project for the given project name in the current directory.
  The command 'new' creates a folder named [projectname] [-qt=false] and generates the following structure:

            ├── CMakeLists.txt
            ├── .clang-format
            ├── run.bat
            ├── {{"src"|foldername}}
            │     └── CMakeLists.txt
            │     └── {{"Mode"|foldername}}
            |          └── mode.cpp
            |          └── mode.h
            ├── {{"test"|foldername}}
            │     └── CMakeLists.txt
            │     └── test.cpp
            ├── {{"vendor"|foldername}}
            │     └── CMakeLists.txt
            ├── {{"docs"|foldername}}
`,
	PreRun: nil,
	Run:    CreateProject,
}

func init() {
	CmdNew.Flag.Var(&qt, "qt", "Support qt application,default false")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdNew)
}

func CreateProject(cmd *commands.Command, args []string) int {
	output := cmd.Out()
	if len(args) == 0 {
		logger.Log.Fatal("Argument [appname] is missing")
	}

	if len(args) > 2 {
		err := cmd.Flag.Parse(args[1:])
		if err != nil {
			logger.Log.Fatal("Parse args err" + err.Error())
		}
	}

	var (
		projectPath string
		// packPath string
		// err      error
	)

	// TODO  添加QT支持
	if qt == "true" {
		logger.Log.Info("Generate new project support GOPATH")
		version.ShowShortVersionBanner()

	} else {
		logger.Log.Info("Generate new project support go mudules.")
		projectPath = path.Join(utils.GetZelWorkPath(), args[0])
		// packPath = args[0]
		if zelVersion.String() == `` {
			zelVersion.Set(utils.ZEL_VERSION)
		}
	}

	if utils.IsExist(projectPath) {
		logger.Log.Errorf(colors.Bold("Application '%s' already exists"), projectPath)
		logger.Log.Warn(colors.Bold("Do you want to overwrite it? [Yes]|No "))
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	logger.Log.Info("Creating C++ project...")
	logger.Log.Infof("project path: %s", projectPath)

	// 创建C++项目所需文件夹
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", projectPath+string(path.Separator), "\x1b[0m")
	os.MkdirAll(projectPath, 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(projectPath, "src")+string(path.Separator), "\x1b[0m")
	os.MkdirAll(path.Join(projectPath, "src"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(projectPath, "src", "Utils")+string(path.Separator), "\x1b[0m")
	os.MkdirAll(path.Join(projectPath, "src", "Utils"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(projectPath, "test")+string(path.Separator), "\x1b[0m")
	os.MkdirAll(path.Join(projectPath, "test"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(projectPath, "vendor")+string(path.Separator), "\x1b[0m")
	os.MkdirAll(path.Join(projectPath, "vendor"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(projectPath, "docs")+string(path.Separator), "\x1b[0m")
	os.MkdirAll(path.Join(projectPath, "docs"), 0755)

	// 创建C++项目所需文件
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(projectPath, "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(path.Join(projectPath, "CMakeLists.txt"), strings.Replace(projectCmakeList, "{{.ProjectName}}", path.Base(args[0]), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(projectPath, ".clang-format"), "\x1b[0m")
	utils.WriteToFile(path.Join(projectPath, ".clang-format"), clangFormat)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(projectPath, "run.bat"), "\x1b[0m")
	utils.WriteToFile(path.Join(projectPath, "run.bat"), runBat)

	return 0
}
