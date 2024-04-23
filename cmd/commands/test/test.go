package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZEL-30/zel/cmake"
	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/logger/colors"
	"github.com/ZEL-30/zel/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var CmdTest = &commands.Command{
	UsageLine: "test [appname] [watchall] [-main=*.go] [-downdoc=true]  [-gendoc=true] [-vendor=true] [-e=folderToExclude] [-ex=extraPackageToWatch] [-tags=goBuildTags] [-runmode=BEEGO_RUNMODE]",
	Short:     "Test the application by starting a local development server",
	Long: `
Run command will supervise the filesystem of the application for any changes, and recompile/restart it.
	`,
	PreRun: func(cmd *commands.Command, args []string) {},
	Run:    RunTest,
}

var (
	new string
)

func init() {
	CmdTest.Flag.StringVar(&new, "new", "", "Create a new test case")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdTest)
}

func RunTest(cmd *commands.Command, args []string) int {
	// if len(args) == 0 {
	// 	logger.Log.Fatal("Argument [testname] is missing")
	// }

	if len(args) > 2 {
		err := cmd.Flag.Parse(args[1:])
		if err != nil {
			logger.Log.Fatal("Parse args err" + err.Error())
		}
	}

	if new != "" {
		newTest(cmd)
	} else {
		runTest(args)
	}

	return 0
}

func runTest(args []string) {
	// 默认应用程序路径是当前工作目录
	appPath, _ := os.Getwd()

	testProgram := args[0] + "-test"

	// logger.Log.Infof("Using '%s' as 'appname'", appName)

	buildPath := filepath.Join(appPath, "build")
	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildType:             config.Conf.BuildType,
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		AppPath:               appPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildType: config.Conf.BuildType,
	}

	testName := cases.Title(language.English).String(args[0])
	fmt.Printf("=== RUN   Test%s\n", testName)
	err := cmake.Run(&configArg, &buildArg, testProgram)
	if err != nil {
		fmt.Printf("--- FAIL: Test%s (0.14s)\n", testName)
		fmt.Println(err.Error())
	} else {
		fmt.Printf("--- PASS: Test%s (0.00s)\n", testName)
	}
}

func newTest(cmd *commands.Command) {

	output := cmd.Out()

	var (
		testPath     string
		testsPath    string
		testFileName string
	)

	testsPath = filepath.Join(utils.GetZelWorkPath(), "tests")
	testPath = filepath.Join(utils.GetZelWorkPath(), "tests", new)
	testFileName = new + "_test.cpp"

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
	utils.WriteToFile(filepath.Join(testPath, testFileName), testContent)

	// 向 cmakelists 中 追加写入内容
	fmt.Fprintf(output, "\t%s%sadd%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filepath.Join(testsPath, "CMakeLists.txt"), "\x1b[0m")
	file, err := os.OpenFile(filepath.Join(testsPath, "CMakeLists.txt"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Log.Fatalf("Open '%s' err: %s", filepath.Join(testsPath, "CMakeLists.txt"), err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(strings.Replace(testCmakeLists, "{{ .TestName }}", filepath.Base(new), -1))
	if err != nil {
		logger.Log.Fatalf("Write '%s' err: %s", filepath.Join(testsPath, "CMakeLists.txt"), err.Error())
	}

	logger.Log.Success("New test case successfully created!")
}
