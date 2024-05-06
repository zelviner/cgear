package install

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ZEL-30/zel/cmake"
	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/cmd/commands/version"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/logger/colors"
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
		logger.Log.Fatal("Please specify third-party library information, for example: zel install google:googletest")
	}

	cmd.Flag.Parse(args[1:])

	// projectPath, _ := os.Getwd()
	vendorInfo := args[0]

	re, err := regexp.Compile("(.+):(.+)")
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	if !re.MatchString(vendorInfo) {
		logger.Log.Fatal("Please specify the correct third-party library information, for example: google:googletest")
	}

	Author := re.FindStringSubmatch(vendorInfo)[1]
	repositoryName := re.FindStringSubmatch(vendorInfo)[2]

	ssh := "git@github.com:" + Author + "/" + repositoryName
	vendorPath = filepath.Join(zelCPath, "pkg", repositoryName)

	if utils.IsExist(vendorPath) {
		logger.Log.Errorf(colors.Bold("%s '%s' already exists"), vendorInfo, vendorPath)
		logger.Log.Warn(colors.Bold("Do you want to update it? [Yes]|No "))
		if utils.AskForConfirmation() {
			logger.Log.Infof("'%s' already exists, updating ...", vendorInfo)
			os.RemoveAll(vendorPath)
		} else {
			logger.Log.Infof("Installing '%s' ...", vendorInfo)
			err = install()
			if err != nil {
				logger.Log.Fatal(err.Error())
			}
			logger.Log.Successf("Successfully installed '%s'", vendorInfo)
			return 0
		}
	}

	err = DownloadPKG(ssh, vendorPath)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Infof("Installing '%s' ...", vendorInfo)
	err = install()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	logger.Log.Successf("Successfully installed '%s'", vendorInfo)
	return 0
}

func DownloadPKG(ssh string, vendorPath string) error {

	logger.Log.Info("Downloading third-party libraries: " + vendorPath)

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

	installPath := filepath.Join(zelCPath, strings.ToLower(buildMode))

	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildMode:             buildMode,
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		AppPath:               vendorPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
		InstallPrefix:         installPath,
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
