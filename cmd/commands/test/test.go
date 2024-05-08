package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ZEL-30/zel/cmake"
	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
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
	new     string // 新建测试用例
	rebuild bool   // 是否重新构建
)

func init() {
	CmdTest.Flag.BoolVar(&rebuild, "r", false, "Clear the build folder in the project and rebuild, default false")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdTest)
}

func RunTest(cmd *commands.Command, args []string) int {
	// if len(args) == 0 {
	// 	logger.Log.Fatal("Argument [testname] is missing")
	// }
	// cmd.Flag.Parse(args[1:])

	if len(args) > 2 {
		err := cmd.Flag.Parse(args[1:])
		if err != nil {
			logger.Log.Fatal("Parse args err" + err.Error())
		}
	}

	// 默认应用程序路径是当前工作目录
	appPath, _ := os.Getwd()

	var (
		testProgram string
		testName    string
	)
	if index := strings.Index(args[0], "."); index == -1 {
		testProgram = args[0] + "-test.exe"
		fmt.Println(testProgram)
	} else {
		testProgram = args[0][:index] + "-test.exe"
		testName = args[0][index+1:]
		fmt.Println(testProgram)
		fmt.Println(testName)
	}

	buildPath := filepath.Join(appPath, "build")
	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildMode:             config.Conf.BuildMode,
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		AppPath:               appPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildMode: config.Conf.BuildMode,
	}

	// testName := cases.Title(language.English).String(args[0])
	err := cmake.Build(&configArg, &buildArg, rebuild, false)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	runPath := filepath.Join(appPath, "build", "test", testProgram)

	arg := fmt.Sprintf("--gtest_filter='%s'", args[0])
	c := exec.Command(runPath, arg)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	return 0
}
