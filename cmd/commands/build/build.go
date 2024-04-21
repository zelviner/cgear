package build

import (
	"os"
	"path/filepath"

	"github.com/ZEL-30/zel/cmake"
	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/cmd/commands/version"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
)

var CmdBuild = &commands.Command{
	UsageLine: "build [-rebuild=false]",
	Short:     "Compile the application",
	Long: `
Build command will supervise the filesystem of the application for any changes, and recompile/restart it.

`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    BuildApp,
}

var (
	rebuild   bool   // 是否重新构建
	appPath   string // 应用程序路径
	buildPath string // 构建路径
)

func init() {
	CmdBuild.Flag.BoolVar(&rebuild, "rebuild", false, "Delete previously compiled build folders and recompile the application.")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdBuild)
}

func BuildApp(cmd *commands.Command, args []string) int {
	appPath, _ = os.Getwd()
	buildPath = filepath.Join(appPath, "build")

	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildType:             0,
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		AppPath:               appPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildType: 0,
	}

	cmake.Build(&configArg, &buildArg, rebuild)

	logger.Log.Success("Build successful!")

	return 0
}
