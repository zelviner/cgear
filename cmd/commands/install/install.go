package install

import (
	"fmt"
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
	UsageLine: "install []",
	Short:     "Downloading and installing C++ third-party open source libraries from GitHub",
	Long: `
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    install,
}

var (
	vendorPath string
	vendorInfo string

	zelHome = os.Getenv("ZEL_HOME")
)

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdInstall)
}

func install(cmd *commands.Command, args []string) int {

	switch len(args) {
	case 0:
		vendorPath = utils.GetZelWorkPath()
		vendorInfo = filepath.Base(vendorPath)
	case 1:
		cmd.Flag.Parse(args[1:])
		vendorInfo = args[0]
		if filepath.IsAbs(vendorInfo) {
			releaseInstall()
			return 0
		} else {
			getPKG(true)
		}
	default:
		logger.Log.Fatal("Too many parameters")
	}

	logger.Log.Infof("Installing '%s' ...", vendorInfo)
	err := compileInstall(true)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Successf("Successfully installed '%s'", vendorInfo)
	return 0
}

func getPKG(showInfo bool) {

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
	vendorPath = filepath.Join(zelHome, "pkg", repositoryName)

	if utils.IsExist(vendorPath) {
		logger.Log.Errorf(colors.Bold("%s '%s' already exists"), vendorInfo, vendorPath)
		logger.Log.Warn(colors.Bold("Do you want to update it? [Yes]|No "))
		if utils.AskForConfirmation() {
			logger.Log.Infof("'%s' already exists, updating ...", vendorInfo)
			os.RemoveAll(vendorPath)
		} else {
			return
		}
	}

	err = downloadPKG(ssh, vendorPath, showInfo)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func downloadPKG(ssh string, vendorPath string, showInfo bool) error {

	logger.Log.Info("Downloading third-party libraries: " + vendorPath)

	command := exec.Command("git", "clone", ssh, vendorPath, "--depth=1")
	if showInfo {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}
	err := command.Run()
	if err != nil {
		return err
	}

	return nil
}

func compileInstall(showInfo bool) error {
	buildPath := filepath.Join(vendorPath, "build")
	buildType := "Debug"
	installPath := filepath.Join(zelHome, strings.ToLower(buildType))

	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildType:             buildType,
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		ProjectPath:           vendorPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
		InstallPrefix:         installPath,
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildType: buildType,
		Target:    "install",
	}

	// Debug
	err := cmake.Build(&configArg, &buildArg, true, showInfo)
	if err != nil {
		return err
	}

	// Release
	buildType = "Release"
	installPath = filepath.Join(zelHome, strings.ToLower(buildType))
	configArg.BuildType = buildType
	configArg.InstallPrefix = installPath
	buildArg.BuildType = buildType
	err = cmake.Build(&configArg, &buildArg, true, showInfo)
	if err != nil {
		return err
	}

	return nil
}

func InstallGTest() {
	gtestPath := filepath.Join(zelHome, "debug", "include", "gtest")
	if utils.IsExist(gtestPath) {
		return
	}

	vendorInfo = "google:googletest"
	getPKG(false)

	cmakePath := vendorPath + "/CMakeLists.txt"
	str := utils.ReadFile(cmakePath)
	os.Remove(cmakePath)
	content := `# For Windows: Prevent overriding the parent project's compiler/linker settings
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)` + "\n\n" + str
	utils.WriteToFile(cmakePath, content)

	err := compileInstall(false)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	// 删除 zel.json 文件
	zelJsonPath := utils.GetZelWorkPath() + "/zel.json"
	if utils.IsExist(zelJsonPath) {
		os.Remove(zelJsonPath)
	}
}

func releaseInstall() {
	logger.Log.Info("Installing third-party libraries: " + vendorInfo)
	logger.Log.Info("Please set the third-party library name:")
	vendorInfo = utils.ReadLine()

	fmt.Println("getZelWorkPath: ", utils.GetZelWorkPath())
}
