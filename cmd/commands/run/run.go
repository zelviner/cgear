package run

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/cmd/commands/version"
	"github.com/ZEL-30/zel/logger"
)

var CmdRun = &commands.Command{
	UsageLine: "run [appname]",
	Short:     "Run the application",
	Long: `
Run command will supervise the filesystem of the application for any changes, and recompile/restart it.

`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    RunApp,
}

var (
	appName  string    // 应用程序名称
	currPath string    // 应用程序路径
	runMode  string    // 当前运行模式
	runArgs  string    // 运行应用程序的额外参数
	exit     chan bool // 发出退出信号的通道
)

func init() {
	CmdRun.Flag.StringVar(&runMode, "runmode", "", "Set the C++ run mode.")
	CmdRun.Flag.StringVar(&runArgs, "runargs", "", "Extra args to run application.")
	exit = make(chan bool)
	commands.AvailableCommands = append(commands.AvailableCommands, CmdRun)
}

// RunApp定位要监视的文件，并启动 C++ 应用程序
func RunApp(cmd *commands.Command, args []string) int {
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

	execmd := exec.Command(".\\run.bat", appName)
	execmd.Stdout = os.Stdout
	execmd.Stderr = os.Stderr
	err := execmd.Run()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	return 0

}
