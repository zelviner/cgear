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
	"github.com/ZEL-30/zel/logger/colors"
	"github.com/ZEL-30/zel/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	zelHome      = utils.GetZelHomePath()
	zelPkg       = utils.GetZelPkgPath()
	zelInstalled = utils.GetZelInstalledPath()
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
	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildType:             buildType,
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		ProjectPath:           vendorPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
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
	configArg.BuildType = buildType
	buildArg.BuildType = buildType

	err = cmake.Build(&configArg, &buildArg, true, showInfo)
	if err != nil {
		return err
	}

	return nil
}

// 自定义列表项
type archItem struct {
	title, desc string
}

func (i archItem) Title() string       { return i.title }
func (i archItem) Description() string { return i.desc }
func (i archItem) FilterValue() string { return i.title }

// 主模型
type archModel struct {
	list   list.Model
	choice string
}

// 初始化模型
func newArchModel() archModel {
	items := []list.Item{
		archItem{
			title: "x86-windows",
			desc:  "32位 Windows 架构（推荐）",
		},
		archItem{
			title: "x64-windows",
			desc:  "64位 Windows 架构",
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 40, 14)
	l.Title = "请选择目标架构"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#3C3C3C")).
		Padding(0, 1)

	return archModel{list: l}
}

func (m archModel) Init() tea.Cmd { return nil }

func (m archModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if selected, ok := m.list.SelectedItem().(archItem); ok {
				m.choice = selected.title
			}
			return m, tea.Quit
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m archModel) View() string {
	return "\n" + m.list.View()
}

// 使用选择器的函数
func selectArch() string {
	p := tea.NewProgram(newArchModel())
	m, err := p.Run()
	if err != nil {
		logger.Log.Fatalf("选择架构失败:", err)
	}

	if model, ok := m.(archModel); ok && model.choice != "" {
		return model.choice
	}
	logger.Log.Fatal("未选择架构")
	return ""
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

	// 调用选择器
	archDir := selectArch()

	// 拷贝 vendorInfo 下的 include 和 lib 目录到 releasePath 下
	releasePath := filepath.Join(zelInstalled, archDir)
	utils.CopyDir(includePath, filepath.Join(releasePath, "include", repositoryName))
	utils.CopyDir(libPath, filepath.Join(releasePath, "lib", repositoryName))

	// 拷贝 vendorInfo 下的 include 和 lib 目录到 debugPath 下
	debugPath := filepath.Join(releasePath, "debug")
	utils.CopyDir(libPath, filepath.Join(debugPath, "lib", repositoryName))

	logger.Log.Successf("Successfully installed '%s'", repositoryName)
}
