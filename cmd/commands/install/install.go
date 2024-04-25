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
)

var CmdInstall = &commands.Command{
	UsageLine: "install",
	Short:     "",
	Long: `
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    install,
}

var (
	ssh string
)

func init() {
	CmdInstall.Flag.StringVar(&ssh, "ssh", "", "")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdInstall)
}

func install(cmd *commands.Command, args []string) int {

	if ssh == "" {
		logger.Log.Fatal("请指定远程库地址")
	}

	re, err := regexp.Compile("/(.+).git")
	if err != nil {
		logger.Log.Error(err.Error())
	}

	repositoryName := re.FindStringSubmatch(ssh)[1]
	zelCPath := os.Getenv("ZEL_C_PATH")
	repositoryPath := filepath.Join(zelCPath, "pkg", repositoryName)
	logger.Log.Info("正在安装远程库: " + repositoryPath)

	// 判断是否存在
	if _, err := os.Stat(repositoryPath); err == nil {
		logger.Log.Info("该远程库已存在")
		os.RemoveAll(repositoryPath)
	}

	command := exec.Command("git", "clone", ssh, repositoryPath, "--depth=1")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Info("远程库安装成功")

	buildPath := filepath.Join(repositoryPath, "build")

	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildMode:             config.Conf.BuildMode,
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		AppPath:               repositoryPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildMode: config.Conf.BuildMode,
	}

	err = cmake.Build(&configArg, &buildArg, true)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Info("远程库编译成功")

	var includePath []string

	// 查找 include 目录
	filepath.Walk(repositoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "include" {
			includePath = append(includePath, path)
			return nil
		}
		return nil
	})

	for _, include := range includePath {
		logger.Log.Info(include)
	}

	// // 复制 include 目录
	// for _, include := range includePath {

	// }

	return 0
}
