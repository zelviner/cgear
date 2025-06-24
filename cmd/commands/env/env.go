package env

import (
	"bytes"
	"fmt"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/env"
	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/utils"
)

const envInfoTemplate string = `%s%s _____     _ 
/ _  / ___| |
\// / / _ \ |
 / //\  __/ |__
/____/\___|___/  v{{ .ZelVersion }}%s
%s%s
├── ZelHome      : {{ .ZelHome }}
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

     $ zel env Toolchain

  ▶ {{"To set BuildType for your C++ project:"|bold}}

     $ zel env BuildType

`,
	Run: SetEnv,
}

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdEnv)
}

func SetEnv(cmd *commands.Command, args []string) int {
	stdout := cmd.Out()
	curryPath := utils.GetZelWorkPath()
	if !utils.IsZelProject(curryPath) {
		logger.Log.Fatal("Not a Zel project")
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
