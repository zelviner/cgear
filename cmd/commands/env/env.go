package env

import (
	"bytes"
	"fmt"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/env"
	"github.com/ZEL-30/zel/logger"
)

const envInfoTemplate string = `%s%s _____     _ 
/ _  / ___| |
\// / / _ \ |
 / //\  __/ |__
/____/\___|___/  v{{ .ZelVersion }}%s
%s%s
├── ZelCPath  : {{ .ZelCPath }}
├── BuildKit  : {{ .BuildKit }}
├── BuildMode : {{ .BuildMode }}
├── TestMode  : {{ .TestMode }}
└── Date      : {{ Now "Monday, 2 Jan 2006" }}%s
`

var CmdEnv = &commands.Command{
	UsageLine: "env [command]",
	Short:     "Setting up the environment for running C++ projects",
	Long: `▶ {{"To set build kit for your C++ project:"|bold}}

     $ zel env kit

  ▶ {{"To set build mode for your C++ project:"|bold}}

     $ zel env mode
  ▶ {{"To set build mode for your C++ project:"|bold}}

     $ zel env mode
`,
	Run: SetEnv,
}

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdEnv)
}

func SetEnv(cmd *commands.Command, args []string) int {
	stdout := cmd.Out()

	if len(args) != 0 {
		gcmd := args[0]
		switch gcmd {
		case "kit":
			env.SetBuildKit()

		case "mode":
			env.SetBuildMode()

		case "test":

		default:
			logger.Log.Fatal("Command is missing")
		}
	} else {
		coloredBanner := fmt.Sprintf(envInfoTemplate, "\x1b[35m", "\x1b[1m",
			"\x1b[0m", "\x1b[32m", "\x1b[1m", "\x1b[0m")
		InitBanner(stdout, bytes.NewBufferString(coloredBanner))
	}

	return 0
}
