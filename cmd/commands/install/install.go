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

// CmdInstall represents the install command
var CmdInstall = &commands.Command{
	UsageLine: "install [package]",
	Short:     "Downloading and installing C++ third-party open source libraries from GitHub",
	Long: `
Install downloads and compiles C++ third-party libraries from GitHub.
Usage:
    zel install                     # Install in current directory
    zel install author:repository   # Install specific repository
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    install,
}

var (
	vendorPath     string
	vendorInfo     string
	repositoryName string

	zelHome   = utils.GetZelHomePath()
	zelPkg    = utils.GetZelPkgPath()
	zelVendor = utils.GetZelVendorPath()
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
		}
		getPKG(true)
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
	repositoryName = re.FindStringSubmatch(vendorInfo)[2]

	ssh := "git@github.com:" + Author + "/" + repositoryName
	vendorPath = filepath.Join(zelPkg, repositoryName)

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

// DownloadPKG downloads a package from GitHub using git clone
// ssh: GitHub SSH URL
// vendorPath: local path to store the package
// showInfo: whether to show download progress
func downloadPKG(ssh string, vendorPath string, showInfo bool) error {

	logger.Log.Info("Downloading third-party libraries: " + repositoryName)

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
	// debug compile
	buildPath := filepath.Join(vendorPath, "build")
	buildType := "Debug"
	installPath := filepath.Join(zelVendor, repositoryName, strings.ToLower(buildType))
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

	err := cmake.Build(&configArg, &buildArg, true, showInfo)
	if err != nil {
		return err
	}

	// release compile
	buildType = "Release"
	installPath = filepath.Join(zelVendor, repositoryName, strings.ToLower(buildType))
	configArg.BuildType = buildType
	configArg.InstallPrefix = installPath
	buildArg.BuildType = buildType

	err = cmake.Build(&configArg, &buildArg, true, showInfo)
	if err != nil {
		return err
	}

	return nil
}

func releaseInstall() {

	// 检测 vendorInfo 是否存在
	if !utils.IsExist(vendorInfo) {
		logger.Log.Fatal("Third-party library not found: " + vendorInfo)
	}

	// 检测 vendorInfo/include 和 vendorInfo/lib 是否存在
	includePath := filepath.Join(vendorInfo, "include")
	libPath := filepath.Join(vendorInfo, "lib")
	if !utils.IsExist(includePath) || !utils.IsExist(libPath) {
		logger.Log.Fatalf("%s is not a third-party library", vendorInfo)
	}

	repositoryName = filepath.Base(vendorInfo)
	logger.Log.Info("Installing third-party libraries: " + vendorInfo)
	logger.Log.Infof("Please set the third-party library name (default: %s):", repositoryName)
	temp := utils.ReadLine()
	if temp != "" {
		repositoryName = temp
	}

	// 拷贝 vendorInfo 下的 include 和 lib 目录到 debugPath 下
	debugPath := filepath.Join(zelVendor, repositoryName, "debug")
	utils.CopyDir(includePath, filepath.Join(debugPath, "include"))
	utils.CopyDir(libPath, filepath.Join(debugPath, "lib"))

	// 拷贝 vendorInfo 下的 include 和 lib 目录到 releasePath 下
	releasePath := filepath.Join(zelVendor, repositoryName, "release")
	utils.CopyDir(includePath, filepath.Join(releasePath, "include"))
	utils.CopyDir(libPath, filepath.Join(releasePath, "lib"))

	logger.Log.Successf("Successfully installed '%s'", repositoryName)
}

func InstallGTest() {
	gtestPath := filepath.Join(zelVendor, "googletest", "debug")
	if utils.IsExist(gtestPath) {
		return
	}

	vendorInfo = "google:googletest"
	getPKG(true)

	cmakePath := vendorPath + "/CMakeLists.txt"
	str := utils.ReadFile(cmakePath)
	os.Remove(cmakePath)
	content := `# For Windows: Prevent overriding the parent project's compiler/linker settings
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)` + "\n\n" + str
	utils.WriteToFile(cmakePath, content)

	err := compileInstall(true)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	// 删除 zel.json 文件
	zelJsonPath := utils.GetZelWorkPath() + "/zel.json"
	if utils.IsExist(zelJsonPath) {
		os.Remove(zelJsonPath)
	}
}
