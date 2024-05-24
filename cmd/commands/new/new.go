package new

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/cmd/commands/version"
	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/logger/colors"
	"github.com/ZEL-30/zel/utils"
)

var (
	test       bool
	qt         bool
	zelVersion utils.DocValue
)

var CmdNew = &commands.Command{
	UsageLine: "new [appname] [-qt=false]",
	Short:     "Create a C++ project, using cmake tool and opening by vscode",
	Long: `Creates a C++ project for the given project name in the current directory.
  The command 'new' creates a folder named [projectname] [-qt=false] and generates the following structure:

            ├── CMakeLists.txt
            ├── .clang-format
            ├── README.md
            ├── {{"src"|foldername}}
            │     └── CMakeLists.txt
            │     └── {{"utils"|foldername}}
            |          └── utils.cpp
            |          └── utils.h
            │     └── main.cpp
            ├── {{"test"|foldername}}
            │     └── CMakeLists.txt
            ├── {{".vecode"|foldername}}
            │     └── launch.json
            ├── {{"docs"|foldername}}
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    Create,
}

func init() {
	CmdNew.Flag.BoolVar(&qt, "qt", false, "New a Qt Application, default false")
	CmdNew.Flag.BoolVar(&test, "test", false, "New a Test Case, default false")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdNew)
}

func Create(cmd *commands.Command, args []string) int {
	if len(args) == 0 {
		logger.Log.Fatal("Argument [appname] is missing")
	}

	err := cmd.Flag.Parse(args[1:])
	if err != nil {
		logger.Log.Fatal("Parse args err" + err.Error())
	}

	switch {
	case qt:
		CreateProjectWithQt(cmd, args)

	case test:
		CreateTestCase(cmd, args[0])

	default:
		CreateProject(cmd, args[0])
	}

	return 0
}

func CreateProject(cmd *commands.Command, appname string) int {

	output := cmd.Out()

	projectPath := filepath.Join(utils.GetZelWorkPath(), appname)

	if utils.IsExist(projectPath) {
		logger.Log.Errorf(colors.Bold("Application '%s' already exists"), projectPath)
		logger.Log.Warn(colors.Bold("Do you want to overwrite it? [Yes]|No "))
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}
	logger.Log.Info("Creating C++ project...")

	// // 下载 GTest 依赖库
	// err := install.DownloadPKG("git@github.com:google/googletest.git", filepath.Join(projectPath, "vendor"))
	// if err != nil {
	// 	logger.Log.Fatal(err.Error())
	// }

	// 创建C++项目所需文件夹
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", projectPath+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(projectPath, 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "src"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src", "utils")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "src", "utils"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "test")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "test"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".vecode")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, ".vscode"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "docs")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "docs"), 0755)

	// 创建C++项目所需文件
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "CMakeLists.txt"), strings.Replace(projectCmakeLists, "{{.ProjectName}}", filepath.Base(appname), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".clang-format"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".clang-format"), clangFormat)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "README.md"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "README.md"), strings.Replace(readme, "{{.ProjectName}}", filepath.Base(appname), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src", "CMakeLists.txt"), srcCmakeLists)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src/utils", "utils.h"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src/utils", "utils.h"), utilsHeader)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src/utils", "utils.cpp"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src/utils", "utils.cpp"), utilsCPP)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "test", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "test", "CMakeLists.txt"), testCmakeLists)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".vsocde", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".vscode", "launch.json"), launch)

	logger.Log.Success("New C++ project successfully created!")
	return 0
}

func CreateProjectWithQt(cmd *commands.Command, args []string) {

}

func CreateTestCase(cmd *commands.Command, testName string) {

	output := cmd.Out()

	var (
		testPath       string
		testsPath      string
		testFileName   string
		testConfigPath string
	)

	testsPath = filepath.Join(utils.GetZelWorkPath(), "test")
	testPath = filepath.Join(utils.GetZelWorkPath(), "test", testName)
	testFileName = testName + "_test.cpp"
	testConfigPath = filepath.Join(utils.GetZelWorkPath(), ".vscode", "launch.json")

	if utils.IsExist(testPath) {
		logger.Log.Errorf(colors.Bold("Test case '%s' already exists"), testPath)
		logger.Log.Warn(colors.Bold("Do you want to overwrite it? [Yes]|No "))
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}
	logger.Log.Info("Creating test case...")

	// 创建C++项目所需文件夹
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", testPath, "\x1b[0m")
	os.MkdirAll(testPath, 0755)

	// 创建C++项目所需文件
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(testPath, testFileName), "\x1b[0m")
	utils.WriteToFile(filepath.Join(testPath, testFileName), strings.Replace(testContent, "{{ .testName }}", testName, -1))
	utils.ReplaceFileContent(testConfigPath, "//{{ .configuration }}", testLaunch)
	utils.ReplaceFileContent(testConfigPath, "{{ .testName }}", testName)

	// 向 test/cmakelists 中 追加写入内容
	fmt.Fprintf(output, "\t%s%sadd%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(testsPath, "CMakeLists.txt"), "\x1b[0m")
	file, err := os.OpenFile(filepath.Join(testsPath, "CMakeLists.txt"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Log.Fatalf("Open '%s' err: %s", filepath.Join(testsPath, "CMakeLists.txt"), err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("add_test_executable(%s)\n", testName))
	if err != nil {
		logger.Log.Fatalf("Write '%s' err: %s", filepath.Join(testsPath, "CMakeLists.txt"), err.Error())
	}

	logger.Log.Success("New test case successfully created!")

}
