package install

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/cmd/commands/version"
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

	projectPath, _ := os.Getwd()

	if ssh == "" {
		logger.Log.Fatal("请指定远程库地址")
	}

	re, err := regexp.Compile("/(.+).git")
	if err != nil {
		logger.Log.Error(err.Error())
	}

	repositoryName := re.FindStringSubmatch(ssh)[1]
	repositoryPath := projectPath + "/vendor/" + repositoryName
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

	return 0
}
