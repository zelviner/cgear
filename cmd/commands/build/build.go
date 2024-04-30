package build

import (
	"os"
	"path/filepath"

	"github.com/ZEL-30/zel/cmake"
	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
)

var CmdBuild = &commands.Command{
	UsageLine: "build [target] [-r]",
	Short:     "Compile the application",
	Long: `
Build command will supervise the filesystem of the application for any changes, and recompile/restart it.

`,
	Run: BuildApp,
}

var (
	rebuild   bool   // 是否重新构建
	target    string // 构建类型
	appPath   string // 应用程序路径
	buildPath string // 构建路径
)

func init() {
	CmdBuild.Flag.BoolVar(&rebuild, "r", false, "Clear the build folder in the project and rebuild, default false")
	CmdBuild.Flag.StringVar(&target, "t", "", "Set the target to compile")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdBuild)
}

func BuildApp(cmd *commands.Command, args []string) int {

	appPath, _ = os.Getwd()
	buildPath = filepath.Join(appPath, "build")

	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildMode:             config.Conf.BuildMode,
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		AppPath:               appPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
		CXXFlags:              "-D_MD",
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildMode: config.Conf.BuildMode,
		Target:    target,
	}

	err := cmake.Build(&configArg, &buildArg, rebuild, true)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Success("Build successful!")

	return 0
}
