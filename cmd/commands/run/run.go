package run

import (
	"path/filepath"

	"github.com/zelviner/cgear/cmake"
	"github.com/zelviner/cgear/cmd/commands"
	"github.com/zelviner/cgear/config"
	"github.com/zelviner/cgear/logger"
	"github.com/zelviner/cgear/utils"
)

var CmdRun = &commands.Command{
	UsageLine: "run [-r]",
	Short:     "Run the application",
	Long: `
Builds the current project and runs the executable named by CMake project().
Use -r to remove the build directory before building.

`,
	PreRun: nil,
	Run:    RunApp,
}

var (
	appName string // 应用程序名称
	rebuild bool   // 是否重建
)

func init() {
	CmdRun.Flag.BoolVar(&rebuild, "r", false, "Clear the build folder in the project and rebuild, default false")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdRun)
}

// RunApp定位要监视的文件，并启动 C++ 应用程序
func RunApp(cmd *commands.Command, args []string) int {

	projectPath := utils.GetCgearWorkPath()
	var err error
	appName, err = utils.GetCgearAppName(projectPath)
	if err != nil {
		logger.Log.Error(err.Error())
		return 1
	}

	buildPath := filepath.Join(projectPath, "build")

	configArg := cmake.ConfigArg{
		Toolchain:             config.Conf.Toolchain,
		Platform:              config.Conf.Platform,
		BuildType:             config.Conf.BuildType,
		Generator:             config.Conf.Generator,
		NoWarnUnusedCli:       true,
		ExportCompileCommands: true,
		ProjectPath:           projectPath,
		BuildPath:             buildPath,
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildType: config.Conf.BuildType,
	}

	if err = cmake.Run(&configArg, &buildArg, appName, rebuild); err != nil {
		return 1
	}

	return 0
}
