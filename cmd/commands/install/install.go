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
	sshArg     string
	vendorPath string

	zelCPath   = os.Getenv("ZEL_C_PATH")
	zelInclude = filepath.Join(zelCPath, "include")
	zelLib     = filepath.Join(zelCPath, "lib")
)

func init() {
	CmdInstall.Flag.StringVar(&sshArg, "ssh", "", "")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdInstall)
}

func installPKG(cmd *commands.Command, args []string) int {

	if sshArg == "" {
		logger.Log.Fatal("请指定远程库地址")
	}

	projectPath, _ := os.Getwd()
	vendorPath = filepath.Join(projectPath, "vendor")
	err := DownloadPKG(sshArg, vendorPath)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	// err = build()
	// if err != nil {
	// 	logger.Log.Fatal(err.Error())
	// }

	// err = install()
	// if err != nil {
	// 	logger.Log.Fatal(err.Error())
	// }

	return 0
}

func DownloadPKG(ssh string, vendorPath string) error {

	re, err := regexp.Compile("/(.+).git")
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	repositoryName := re.FindStringSubmatch(ssh)[1]
	vendorPath = filepath.Join(vendorPath, repositoryName)

	logger.Log.Info("正在下载远程库: " + vendorPath)

	// 判断是否存在
	if utils.FileIsExisted(vendorPath) {
		logger.Log.Info("该远程库已存在")
		os.RemoveAll(vendorPath)
	}

	command := exec.Command("git", "clone", ssh, vendorPath, "--depth=1")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		return err
	}

	return nil
}

func build() error {
	buildPath := filepath.Join(vendorPath, "build")

	logger.Log.Info("正在编译远程库: " + buildPath)

	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildMode:             config.Conf.BuildMode,
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		AppPath:               vendorPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildMode: config.Conf.BuildMode,
	}

	err := cmake.Build(&configArg, &buildArg, true, true)
	if err != nil {
		return err
	}

	logger.Log.Info("远程库编译成功")
	return nil
}

func install() error {

	// var includePaths []string
	// var libPath []string

	// 查找 include 目录
	err := filepath.Walk(vendorPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "include" {
			logger.Log.Info("找到 include 目录: " + path)
			err = utils.CopyDir(path, zelInclude)
			return err
		}

		if info.IsDir() && info.Name() == "lib" {
			logger.Log.Info("找到 lib 目录: " + path)
			err = utils.CopyDir(path, zelLib)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// // 复制 include 目录
	// for _, include := range includePath {

	// }

	return nil
}
