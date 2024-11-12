package pack

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

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
)

func init() {
	// fs := flag.NewFlagSet("pack", flag.ContinueOnError)
	// fs.StringVar(&projectPath, "p", "", "Set the project path. Defaults to the current path.")
	// fs.BoolVar(&isBuild, "b", true, "Tell the command to do a build for the current platform. Defaults to true.")
	// fs.StringVar(&projectName, "a", "", "Set the application name. Defaults to the dir name.")
	// fs.StringVar(&buildArgs, "ba", "", "Specify additional args for Go build.")
	// // fs.Var(&buildEnvs, "be", "Specify additional env variables for Go build. e.g. GOARCH=arm.")
	// fs.StringVar(&outputP, "o", "", "Set the compressed file output path. Defaults to the current path.")
	// fs.StringVar(&format, "f", "tar.gz", "Set file format. Either tar.gz or zip. Defaults to tar.gz.")
	// fs.StringVar(&excludeP, "exp", ".", "Set prefixes of paths to be excluded. Uses a column (:) as separator.")
	// fs.StringVar(&excludeS, "exs", ".go:.DS_Store:.tmp", "Set suffixes of paths to be excluded. Uses a column (:) as separator.")
	// // fs.Var(&excludeR, "exr", "Set a regular expression of files to be excluded.")
	// fs.BoolVar(&fsym, "fs", false, "Tell the command to follow symlinks. Defaults to false.")
	// fs.BoolVar(&ssym, "ss", false, "Tell the command to skip symlinks. Defaults to false.")
	// fs.BoolVar(&verbose, "v", false, "Be more verbose during the operation. Defaults to false.")
	// CmdPack.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdPack)
}

func packProject(cmd *commands.Command, args []string) int {
	projectPath = utils.GetZelWorkPath()
	projectName := filepath.Base(projectPath)

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

	logger.Log.Infof("Packaging Project on '%s'...", projectPath)

	var (
		versionNumber string
	)
	logger.Log.Infof("Please set the version number: ")
	fmt.Scanf("%s", &versionNumber)

	// 编译
	build()

	desPath := filepath.Join(projectPath, "bin", projectName+"-"+versionNumber)
	des := filepath.Join(desPath, projectName+"-"+versionNumber+".exe")
	zipdir := desPath + ".zip"

	utils.MakeDir(desPath)
	defer func() {
		// Remove the desPath once bee pack is done
		err := os.RemoveAll(desPath)
		if err != nil {
			logger.Log.Error("Failed to remove the generated des dir")
		}
	}()

	err := filepath.Walk(filepath.Join(projectPath, "bin"), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".exe" || filepath.Ext(path) == ".dll" {
			if info.Name() == projectName+".exe" {
				_, err := utils.CopyFile(path, des)
				if err != nil {
					return err
				}
			} else {
				_, err := utils.CopyFile(path, filepath.Join(desPath, info.Name()))
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	// QT 打包
	err = pack(des)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	// 压缩
	err = utils.ZipFile(desPath, zipdir)
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
