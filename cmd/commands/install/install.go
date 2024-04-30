package install

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/ZEL-30/zel/cmake"
	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/cmd/commands/version"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/utils"
)

var CmdInstall = &commands.Command{
	UsageLine: "install",
	Short:     "",
	Long: `
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    installPKG,
}

var (
	vendorPath string
	isDebug    bool

	zelCPath = os.Getenv("ZEL_C_PATH")
)

func init() {
	CmdInstall.Flag.BoolVar(&isDebug, "d", false, "编译Release模式")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdInstall)
}

func installPKG(cmd *commands.Command, args []string) int {

	if len(args) < 1 {
		logger.Log.Fatal("请指定第三方库信息, 例如: zel install ZEL-30:zel")
	}

	cmd.Flag.Parse(args[1:])

	// projectPath, _ := os.Getwd()
	vendorInfo := args[0]

	re, err := regexp.Compile("(.+):(.+)")
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	if !re.MatchString(vendorInfo) {
		logger.Log.Fatal("请指定正确的第三方库信息, 例如: google:googletest")
	}

	Author := re.FindStringSubmatch(vendorInfo)[1]
	repositoryName := re.FindStringSubmatch(vendorInfo)[2]

	ssh := "git@github.com:" + Author + "/" + repositoryName
	vendorPath = filepath.Join(zelCPath, "pkg", repositoryName)

	// 判断是否存在
	if utils.FileIsExisted(vendorPath) {
		logger.Log.Infof("'%s' 已存在, 更新中...", vendorInfo)
		os.RemoveAll(vendorPath)
	}

	err = DownloadPKG(ssh, vendorPath)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Infof("正在安装 '%s'", vendorInfo)
	err = install()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	logger.Log.Successf("'%s' 安装成功", vendorInfo)
	return 0
}

func DownloadPKG(ssh string, vendorPath string) error {

	logger.Log.Info("正在下载远程库: " + vendorPath)

	command := exec.Command("git", "clone", ssh, vendorPath, "--depth=1")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		return err
	}

	return nil
}

func install() error {
	buildPath := filepath.Join(vendorPath, "build")

	buildMode := "Release"
	if isDebug {
		buildMode = "Debug"
	}

	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildMode:             buildMode,
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		AppPath:               vendorPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
		InstallPrefix:         zelCPath,
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildMode: buildMode,
		Target:    "install",
	}

	err := cmake.Build(&configArg, &buildArg, true, true)
	if err != nil {
		return err
	}

	return nil
}
