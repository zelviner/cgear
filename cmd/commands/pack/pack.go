package pack

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	path "path/filepath"
	"strings"

	"github.com/ZEL-30/zel/cmake"
	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/utils"
)

var CmdPack = &commands.Command{
	UsageLine: "pack",
	Short:     "Compresses a C++ project into a single file",
	Long: `Pack is used to compress C++ project into a tarball/zip file.
  This eases the deployment by directly extracting the file to a server.
 
  {{"Example:"|bold}}
    $ bee pack -v -ba="-ldflags '-s -w'" 
`,
	PreRun: func(cmd *commands.Command, args []string) {},
	Run:    packProject,
}

var (
	projectPath string
	projectName string
	excludeP    string
	excludeS    string
	outputP     string
	// excludeR  utils.ListOpts
	fsym      bool
	ssym      bool
	isBuild   bool
	buildArgs string
	// buildEnvs utils.ListOpts
	verbose bool
	format  string
)

func init() {
	fs := flag.NewFlagSet("pack", flag.ContinueOnError)
	fs.StringVar(&projectPath, "p", "", "Set the project path. Defaults to the current path.")
	fs.BoolVar(&isBuild, "b", true, "Tell the command to do a build for the current platform. Defaults to true.")
	fs.StringVar(&projectName, "a", "", "Set the application name. Defaults to the dir name.")
	fs.StringVar(&buildArgs, "ba", "", "Specify additional args for Go build.")
	// fs.Var(&buildEnvs, "be", "Specify additional env variables for Go build. e.g. GOARCH=arm.")
	fs.StringVar(&outputP, "o", "", "Set the compressed file output path. Defaults to the current path.")
	fs.StringVar(&format, "f", "tar.gz", "Set file format. Either tar.gz or zip. Defaults to tar.gz.")
	fs.StringVar(&excludeP, "exp", ".", "Set prefixes of paths to be excluded. Uses a column (:) as separator.")
	fs.StringVar(&excludeS, "exs", ".go:.DS_Store:.tmp", "Set suffixes of paths to be excluded. Uses a column (:) as separator.")
	// fs.Var(&excludeR, "exr", "Set a regular expression of files to be excluded.")
	fs.BoolVar(&fsym, "fs", false, "Tell the command to follow symlinks. Defaults to false.")
	fs.BoolVar(&ssym, "ss", false, "Tell the command to skip symlinks. Defaults to false.")
	fs.BoolVar(&verbose, "v", false, "Be more verbose during the operation. Defaults to false.")
	CmdPack.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdPack)
}

func packProject(cmd *commands.Command, args []string) int {
	currPath := utils.GetZelWorkPath()
	var thePath string

	nArgs := []string{}
	has := false
	for _, a := range args {
		if a != "" && a[0] == '-' {
			has = true
		}
		if has {
			nArgs = append(nArgs, a)
		}
	}
	cmd.Flag.Parse(nArgs)

	if !path.IsAbs(projectPath) {
		projectPath = path.Join(currPath, projectPath)
	}

	thePath, err := path.Abs(projectPath)
	if err != nil {
		logger.Log.Fatalf("Wrong project path: %s", projectPath)
	}

	if stat, err := os.Stat(thePath); os.IsNotExist(err) || !stat.IsDir() {
		logger.Log.Fatalf("Project path does not exist: %s", thePath)
	}

	logger.Log.Infof("Packaging Project on '%s'...", thePath)

	var (
		versionNumber string
		execName      string
	)
	logger.Log.Infof("Please set the version number: ")
	fmt.Scanf("%s", &versionNumber)

	filepath.Walk(filepath.Join(currPath, "bin"), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if index := strings.Index(info.Name(), ".exe"); index != -1 {
			execName = info.Name()[:index]
		}

		return nil
	})

	src := filepath.Join(currPath, "bin", execName+".exe")
	tmpdir := filepath.Join(currPath, "bin", execName+"-"+versionNumber)
	zipdir := tmpdir + ".zip"
	des := filepath.Join(tmpdir, execName+"-"+versionNumber+".exe")

	utils.MakeDir(tmpdir)
	defer func() {
		// Remove the tmpdir once bee pack is done
		err := os.RemoveAll(tmpdir)
		if err != nil {
			logger.Log.Error("Failed to remove the generated temp dir")
		}
	}()

	// 编译
	if isBuild {
		build()
	}

	_, err = utils.CopyFile(src, des)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	switch format {
	case "tar.gz":
	default:
		format = "zip"
	}

	// QT 打包
	err = pack(des)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	// 压缩
	err = utils.ZipFile(tmpdir, zipdir)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Success("Application packed!")
	return 0
}

func pack(execPath string) error {
	cmd := exec.Command("windeployqt", execPath)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func build() {
	buildPath := filepath.Join(projectPath, "build")

	configArg := cmake.ConfigArg{
		NoWarnUnusedCli:       true,
		BuildType:             "Release",
		ExportCompileCommands: true,
		Kit:                   config.Conf.Kit,
		ProjectPath:           projectPath,
		BuildPath:             buildPath,
		Generator:             "Ninja",
		CXXFlags:              "-D_MD",
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildType: "Release",
	}

	err := cmake.Build(&configArg, &buildArg, true, false)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Success("Build successful!")
}
