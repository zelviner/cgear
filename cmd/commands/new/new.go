package new

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/cmd/commands/install"
	"github.com/ZEL-30/zel/cmd/commands/version"
	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/logger/colors"
	"github.com/ZEL-30/zel/utils"
)

var (
	test         bool
	qt           bool
	zelVersion   utils.DocValue
	output       io.Writer
	projectPath  string
	projectName  string
	projectTypes = []string{"Application", "QT Application", "Static library", "Dynamic library", "Test cases"}
)

var CmdNew = &commands.Command{
	UsageLine: "new [project_name] [-qt=false]",
	Short:     "Create a C++ project, using cmake tool and opening by vscode",
	Long: `Creates a C++ project for the given project name in the current directory.
  The command 'new' creates a folder named [projectname] [-qt=false] and generates the following structure:

            ├── CMakeLists.txt
            ├── .clang-format
            ├── .gitignore
            ├── README.md
            ├── LICENSE.txt
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
		logger.Log.Fatal("Argument [projectName] is missing")
	}

	err := cmd.Flag.Parse(args[1:])
	if err != nil {
		logger.Log.Fatal("Parse args err" + err.Error())
	}

	install.InstallGTest()

	output = cmd.Out()
	projectName = args[0]
	projectPath = filepath.Join(utils.GetZelWorkPath(), projectName)

	if utils.IsExist(projectPath) {
		logger.Log.Errorf(colors.Bold("Project '%s' already exists"), projectPath)
		logger.Log.Warn(colors.Bold("Do you want to overwrite it? [Yes]|No "))
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	logger.Log.Info("Please selecte project type:")
	for i, projectType := range projectTypes {
		fmt.Printf("\t[%d] %s\n", i+1, projectType)
	}

	var index int
	fmt.Scanf("%d", &index)

	switch index {
	case 1:
		createApp()

	case 2:
		createQTApp()

	case 3:
		createLib("static library")

	case 4:
		createLib("dynamic library")

	case 5:
		createTestCase()

	default:
		createApp()
	}

	return 0
}

func createApp() {
	logger.Log.Info("Creating application ...")

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
	utils.WriteToFile(filepath.Join(projectPath, "CMakeLists.txt"), strings.Replace(projectCMakeLists, "{{ .ProjectName }}", filepath.Base(projectName), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".clang-format"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".clang-format"), clangFormat)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".gitignore"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".gitignore"), gitignore)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "LICENSE.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "LICENSE.txt"), license)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "README.md"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "README.md"), strings.Replace(readme, "{{ .ProjectName }}", filepath.Base(projectName), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src", "CMakeLists.txt"), strings.Replace(appSrcCMakeLists, "{{ .ProjectName }}", filepath.Base(projectName), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src", "main.cpp"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src", "main.cpp"), mainCPP)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src/utils", "utils.h"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src/utils", "utils.h"), utilsHeader)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src/utils", "utils.cpp"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src/utils", "utils.cpp"), utilsCPP)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "test", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "test", "CMakeLists.txt"), appTestCMakeLists)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".vsocde", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".vscode", "launch.json"), launch)

	logger.Log.Success("New application successfully created!")
}

func createQTApp() {
	logger.Log.Info("Creating qt application ...")

	// 创建C++项目所需文件夹
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", projectPath+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(projectPath, 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "src"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src", "utils")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "src", "utils"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src", "app")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "src", "app"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "test")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "test"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".vecode")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, ".vscode"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "res")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "res"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "res", "images")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "res", "images"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "res", "rc")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "res", "rc"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "res", "ui")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "res", "ui"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "docs")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "docs"), 0755)

	// 创建C++项目所需文件
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "CMakeLists.txt"), strings.Replace(projectCMakeLists, "{{ .ProjectName }}", filepath.Base(projectName), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".clang-format"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".clang-format"), clangFormat)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".gitignore"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".gitignore"), gitignore)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "LICENSE.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "LICENSE.txt"), license)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "README.md"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "README.md"), strings.Replace(readme, "{{ .ProjectName }}", filepath.Base(projectName), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src", "CMakeLists.txt"), strings.Replace(qtSrcCMakeLists, "{{ .ProjectName }}", filepath.Base(projectName), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src", "main.cpp"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src", "main.cpp"), qtMainCPP)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src/app", "main_window.h"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src/app", "main_window.h"), appHeader)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src/app", "main_window.cpp"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src/app", "main_window.cpp"), appCPP)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src/utils", "utils.h"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src/utils", "utils.h"), utilsHeader)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src/utils", "utils.cpp"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src/utils", "utils.cpp"), utilsCPP)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "test", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "test", "CMakeLists.txt"), appTestCMakeLists)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".vsocde", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".vscode", "launch.json"), launch)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".vsocde", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".vscode", "launch.json"), launch)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "res/rc", "logo.rc"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "res/rc", "logo.rc"), logoRc)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "res/rc", "images.qrc"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "res/rc", "images.qrc"), imagesRC)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "res/ui", "main_window.ui"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "res/ui", "main_window.ui"), mainWindowUI)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "res/ui", "template.ui"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "res/ui", "template.ui"), templateUI)

	logger.Log.Success("New qt application successfully created!")
}

func createLib(libType string) {
	logger.Log.Infof("Creating %s ...", libType)

	libSrcCMakeLists = strings.Replace(libSrcCMakeLists, "{{ .ProjectName }}", filepath.Base(projectName), -1)
	if libType == "static library" {
		libSrcCMakeLists = strings.Replace(libSrcCMakeLists, "{{ .LibInfo }}", filepath.Base(staticLibInfo), -1)
	} else {
		libSrcCMakeLists = strings.Replace(libSrcCMakeLists, "{{ .LibInfo }}", filepath.Base(dynamicLibInfo), -1)
	}

	// 创建C++项目所需文件夹
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", projectPath+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(projectPath, 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "include")+string(filepath.Separator), "\x1b[0m")
	os.MkdirAll(filepath.Join(projectPath, "include"), 0755)
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
	utils.WriteToFile(filepath.Join(projectPath, "CMakeLists.txt"), strings.Replace(projectCMakeLists, "{{ .ProjectName }}", filepath.Base(projectName), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".clang-format"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".clang-format"), clangFormat)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".gitignore"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".gitignore"), gitignore)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "LICENSE.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "LICENSE.txt"), license)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "README.md"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "README.md"), strings.Replace(readme, "{{ .ProjectName }}", filepath.Base(projectName), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "include", projectName+".h"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "include", projectName+".h"), "#pragma once\n\n #include <utils/utils.h>")
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src", "CMakeLists.txt"), libSrcCMakeLists)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src/utils", "utils.h"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src/utils", "utils.h"), utilsHeader)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "src/utils", "utils.cpp"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "src/utils", "utils.cpp"), utilsCPP)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, "test", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, "test", "CMakeLists.txt"), testCMakeLists)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(projectPath, ".vsocde", "CMakeLists.txt"), "\x1b[0m")
	utils.WriteToFile(filepath.Join(projectPath, ".vscode", "launch.json"), launch)

	logger.Log.Successf("New %s successfully created!", libType)
}

func createTestCase() {

	var (
		testPath       string
		testsPath      string
		testFileName   string
		testConfigPath string
	)

	testsPath = filepath.Join(utils.GetZelWorkPath(), "test")
	testPath = filepath.Join(utils.GetZelWorkPath(), "test", projectName)
	testFileName = projectName + "_test.cpp"
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
	utils.WriteToFile(filepath.Join(testPath, testFileName), strings.Replace(testContent, "{{ .testName }}", utils.CapitalizeFirstLetter(projectName), -1))
	utils.ReplaceFileContent(testConfigPath, "//{{ .configuration }}", testLaunch)
	utils.ReplaceFileContent(testConfigPath, "{{ .testName }}", projectName)

	// 向 test/cmakelists 中 追加写入内容
	fmt.Fprintf(output, "\t%s%sadd%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(testsPath, "CMakeLists.txt"), "\x1b[0m")
	file, err := os.OpenFile(filepath.Join(testsPath, "CMakeLists.txt"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Log.Fatalf("Open '%s' err: %s", filepath.Join(testsPath, "CMakeLists.txt"), err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("add_test_executable(%s)\n", projectName))
	if err != nil {
		logger.Log.Fatalf("Write '%s' err: %s", filepath.Join(testsPath, "CMakeLists.txt"), err.Error())
	}

	logger.Log.Success("New test case successfully created!")
}
