package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/cmd/commands/version"
	"github.com/ZEL-30/zel/logger"
)

var CmdBuild = &commands.Command{
	UsageLine: "build",
	Short:     "Compile the application",
	Long: `
Build command will supervise the filesystem of the application for any changes, and recompile/restart it.

`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    BuildApp,
}

var (
	appName  string // 应用程序名称
	currPath string // 应用程序路径
)

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdBuild)
}

func BuildApp(cmd *commands.Command, args []string) int {
	// 默认应用程序路径是当前工作目录
	appPath, _ := os.Getwd()

	// 如果提供了参数，我们将其用作应用程序路径
	if len(args) != 0 && args[0] != "watchall" {
		if filepath.IsAbs(args[0]) {
			appPath = args[0]
		} else {
			appPath = filepath.Join(appPath, args[0])
		}
	}

	appName = filepath.Base(appPath)

	logger.Log.Infof("Using '%s' as 'appname'", appName)

	args = append([]string{appName}, args...)
	build(args)

	return 0
}

func build(args []string) {

	if len(args) == 0 {
		logger.Log.Fatal("Please provide the executable file name as an argument. For example: .\run.bat hello")
	}

	var (
		buildDir       = "build"
		executablePath = "bin"
		executableName = fmt.Sprintf("%s.exe", args[0])
	// vsPath
	)

	logger.Log.Infof("%s, %s, %s", buildDir, executablePath, executableName)

	// works, err := os.Stat("C:/Workspaces")
	// if err != nil {
	// 	logger.Log.Error("file not found")
	// }

	// execmd := exec.Command(".\\build.bat", appName)
	// execmd.Stdout = os.Stdout
	// execmd.Stderr = os.Stderr
	// err := execmd.Run()
	// if err != nil {
	// 	logger.Log.Fatal(err.Error())
	// }

}
