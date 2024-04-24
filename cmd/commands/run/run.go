package run

import (
	"os"
	"path/filepath"

	"github.com/ZEL-30/zel/cmake"
	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/config"
)

var CmdRun = &commands.Command{
	UsageLine: "run [appname]",
	Short:     "Run the application",
	Long: `
Run command will supervise the filesystem of the application for any changes, and recompile/restart it.

`,
	PreRun: nil,
	Run:    RunApp,
}

var (
	appName  string    // 应用程序名称
	currPath string    // 应用程序路径
	rebuild  bool      // 是否重建
	exit     chan bool // 发出退出信号的通道
)

func init() {
	CmdRun.Flag.BoolVar(&rebuild, "r", false, "Clear the build folder in the project and rebuild, default false")
	exit = make(chan bool)
	commands.AvailableCommands = append(commands.AvailableCommands, CmdRun)
}

// RunApp定位要监视的文件，并启动 C++ 应用程序
func RunApp(cmd *commands.Command, args []string) int {
	cmd.Flag.Parse(args[1:])

	// 默认应用程序路径是当前工作目录
	appPath, _ := os.Getwd()

	// // 如果提供了参数，我们将其用作应用程序路径
	// if len(args) != 0 && args[0] != "watchall" {
	// 	if filepath.IsAbs(args[0]) {
	// 		appPath = args[0]
	// 	} else {
	// 		appPath = filepath.Join(appPath, args[0])
	// 	}
	// }

	appName = filepath.Base(appPath)
	// logger.Log.Infof("Using '%s' as 'appname'", appName)

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

	cmake.Run(&configArg, &buildArg, args[0], rebuild)

	return 0

}
