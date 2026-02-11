package pack

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zelviner/cgear/cmake"
	"github.com/zelviner/cgear/cmd/commands"
	"github.com/zelviner/cgear/config"
	"github.com/zelviner/cgear/logger"
	"github.com/zelviner/cgear/utils"
)

var CmdPack = &commands.Command{
	UsageLine: "pack",
	Short:     "Compresses a C++ project into a single file",
	Long: `Pack is used to compress C++ project into a tarball/zip file.
  This eases the deployment by directly extracting the file to a server.
 
  {{"Example:"|bold}}
    $ cgear pack -v -ba="-ldflags '-s -w'" 
`,
	PreRun: func(cmd *commands.Command, args []string) {},
	Run:    packProject,
}

var (
	projectPath string
)

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdPack)
}

func packProject(cmd *commands.Command, args []string) int {
	projectPath = utils.GetCgearWorkPath()
	projectName, _ := utils.GetCgearAppName(projectPath)

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

	desPath := filepath.Join(projectPath, "bin", "release", projectName+"-"+versionNumber)
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
			} else if strings.Contains(info.Name(), "_test.exe") {
				// Skip the _test.exe file
				return nil
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

	// 运行时依赖
	if err = runtimeDependencies(desPath); err != nil {
		logger.Log.Fatal(err.Error())
	}

	// TODO Innosetup 打包

	// 便携式打包
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
		Toolchain:             config.Conf.Toolchain,
		Platform:              config.Conf.Platform,
		BuildType:             "Release",
		Generator:             config.Conf.Generator,
		NoWarnUnusedCli:       true,
		ExportCompileCommands: true,
		ProjectPath:           projectPath,
		BuildPath:             buildPath,
		CXXFlags:              "-D_MD",
	}

	buildArg := cmake.BuildArg{
		BuildPath: buildPath,
		BuildType: config.Conf.BuildType,
		IsMSVC:    config.Conf.Toolchain.IsMSVC,
	}

	err := cmake.Build(&configArg, &buildArg, true, false)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Success("Build successful!")
}

func runtimeDependencies(desPath string) error {
	var dllPath string
	cgearHome := utils.GetCgearInstalledPath()
	switch config.Conf.Platform {
	case "x86":
		dllPath = filepath.Join(cgearHome, "x86-windows")
	case "x64":
		dllPath = filepath.Join(cgearHome, "x64-windows")
	}

	// switch config.Conf.BuildType {
	// case "Debug":
	// 	dllPath = filepath.Join(dllPath, "debug", "bin")
	// case "Release":
	// 	dllPath = filepath.Join(dllPath, "bin")
	// }

	dllPath = filepath.Join(dllPath, "bin")
	for _, dep := range config.Conf.RuntimeDependencies {
		if dep == "input dynamic libraries here" {
			continue
		}

		dep += ".dll"
		dll := filepath.Join(dllPath, dep)
		if !utils.IsExist(dll) {
			return fmt.Errorf("runtime dependency %s not found in %s", dep, dllPath)
		}

		_, err := utils.CopyFile(dll, filepath.Join(desPath, dep))
		if err != nil {
			return err
		}
	}

	return nil
}

// func innosetupPack() error {
// issPath := filepath.Join(projectPath, "installer.iss")

// // 调用 ISCC.exe
// iscc := `C:\Program Files (x86)\Inno Setup 6\ISCC.exe` // 或者你的安装路径
// cmdISCC := exec.Command(iscc, issPath)
// cmdISCC.Stdout = os.Stdout
// cmdISCC.Stderr = os.Stderr

// logger.Log.Infof("Running Inno Setup compiler...")

// if err := cmdISCC.Run(); err != nil {
// 	logger.Log.Fatal("Inno Setup failed: ", err.Error())
// }

// logger.Log.Success("Inno Setup installer generated successfully!")

// 	return nil
// }
