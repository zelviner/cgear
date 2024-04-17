package test

import (
	"os"
	"os/exec"
	path "path/filepath"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/logger"
)

var CmdTest = &commands.Command{
	UsageLine: "test [appname] [watchall] [-main=*.go] [-downdoc=true]  [-gendoc=true] [-vendor=true] [-e=folderToExclude] [-ex=extraPackageToWatch] [-tags=goBuildTags] [-runmode=BEEGO_RUNMODE]",
	Short:     "Test the application by starting a local development server",
	Long: `
Run command will supervise the filesystem of the application for any changes, and recompile/restart it.
	`,
	PreRun: func(cmd *commands.Command, args []string) {},
	Run:    RunProject,
}

var (
	projectPath string
	currPath    string
	projectName string
)

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdTest)
}

func RunProject(cmd *commands.Command, args []string) int {
	currPath, _ = os.Getwd()

	if !path.IsAbs(projectPath) {
		projectPath = path.Join(currPath, projectPath)
	}

	thePath, err := path.Abs(projectPath)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	if len(projectName) == 0 {
		projectName = path.Base(thePath)
	}

	execmd := exec.Command(".\\run.bat", args...)
	execmd.Stdout = os.Stdout
	execmd.Stderr = os.Stderr
	err = execmd.Run()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Success("Build Successful!")

	return 0
}
