package env

import (
	"bytes"
	"fmt"

	"github.com/zelviner/cgear/cmd/commands"
	"github.com/zelviner/cgear/config"
	"github.com/zelviner/cgear/env"
	"github.com/zelviner/cgear/logger"
	"github.com/zelviner/cgear/utils"
)

const envInfoTemplate string = `
%s%s  ██████╗ ██████╗ ███████╗ █████╗ ██████╗ 
██╔════╝██╔════╝ ██╔════╝██╔══██╗██╔══██╗
██║     ██║  ███╗█████╗  ███████║██████╔╝
██║     ██║   ██║██╔══╝  ██╔══██║██╔══██╗
╚██████╗╚██████╔╝███████╗██║  ██║██║  ██║
 ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝  v{{ .CgearVersion }}%s
%s%s
├── CgearHome    : {{ .CgearHome }}
├── Toolchain    : {{ .Toolchain }}
├── Platform     : {{ .Platform }}
├── Generator    : {{ .Generator }}
├── BuildType    : {{ .BuildType }}
├── ProjectType  : {{ .ProjectType }}
└── Date         : {{ Now "Monday, 2 Jan 2006" }}%s
`

var CmdEnv = &commands.Command{
	UsageLine: "env [command]",
	Short:     "Setting up the environment for running C++ projects",
	Long: `▶ {{"To set Toolchain for your C++ project:"|bold}}

     $ cgear env Toolchain

  ▶ {{"To set BuildType for your C++ project:"|bold}}

     $ cgear env BuildType

`,
	Run: SetEnv,
}

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdEnv)
}

func SetEnv(cmd *commands.Command, args []string) int {
	stdout := cmd.Out()
	curryPath := utils.GetCgearWorkPath()
	if !utils.IsCgearProject(curryPath) {
		logger.Log.Fatal("Not a Cgear project")
	}

	if len(args) != 0 {
		gcmd := args[0]
		switch gcmd {

		case "Toolchain":
			env.SetToolchain()

		case "Generator":
			env.SetGenerator()

		case "Platform":
			env.SetPlatform()

		case "BuildType":
			env.SetBuildType()

		case "test":

		default:
			logger.Log.Fatal("Command is missing")
		}
	} else {
		coloredBanner := fmt.Sprintf(envInfoTemplate, "\x1b[35m", "\x1b[1m",
			"\x1b[0m", "\x1b[32m", "\x1b[1m", "\x1b[0m")
		InitBanner(stdout, bytes.NewBufferString(coloredBanner))
	}

	config.SaveConfig(config.Conf.ProjectPath)
	return 0
}
