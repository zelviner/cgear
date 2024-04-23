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
├── Kit         : {{ .Kit }}
├── BuildType   : {{ .BuildType }}
└── Date        : {{ Now "Monday, 2 Jan 2006" }}%s
`

var CmdEnv = &commands.Command{
	UsageLine: "env [command]",
	Short:     "Setting up the environment for running C++ projects.",
	Long: `▶ {{"To set kit for your C++ project:"|bold}}

     $ zel env kit

  ▶ {{"To set build type for your C++ project:"|bold}}

     $ zel env buildtype
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
			env.SetKit(cmd, args)
		case "buildtype":
			// buildtype(cmd, args, currPath)

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
