package build

import (
	"path/filepath"

	"github.com/zelviner/cgear/cmake"
	"github.com/zelviner/cgear/cmd/commands"
	"github.com/zelviner/cgear/config"
	"github.com/zelviner/cgear/logger"
	"github.com/zelviner/cgear/utils"
)

var CmdBuild = &commands.Command{
	UsageLine: "build [target] [-r]",
	Short:     "Compile the application",
	Long: `
Configures the current Cgear project with CMake and builds it.
Use -t to build a specific target and -r to remove the build directory first.

`,
	Run: BuildApp,
}

var (
	rebuild   bool   // 是否重新构建
	target    string // 构建类型
	buildPath string // 构建路径
)

func init() {
	CmdBuild.Flag.BoolVar(&rebuild, "r", false, "Clear the build folder in the project and rebuild, default false")
	CmdBuild.Flag.StringVar(&target, "t", "", "Set the target to compile")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdBuild)
}

func BuildApp(cmd *commands.Command, args []string) int {

	appPath := utils.GetCgearWorkPath()
	buildPath = filepath.Join(appPath, "build")

	configArg := cmake.ConfigArg{
		Toolchain:             config.Conf.Toolchain,
		Platform:              config.Conf.Platform,
		BuildType:             config.Conf.BuildType,
		Generator:             config.Conf.Generator,
		NoWarnUnusedCli:       true,
		ExportCompileCommands: true,
		ProjectPath:           appPath,
		BuildPath:             buildPath,
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		Target:    target,
		BuildType: config.Conf.BuildType,
	}

	err := cmake.Build(&configArg, &buildArg, rebuild, true)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Success("Build successful!")

	return 0
}
