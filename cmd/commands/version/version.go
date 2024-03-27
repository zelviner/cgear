package version

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"runtime"

	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/logger"
	"github.com/ZEL-30/zel/logger/colors"
)

const verboseVersionBanner string = `%s%s _____     _ 
/ _  / ___| |
\// / / _ \ |
 / //\  __/ |__
/____/\___|___/  v{{ .ZelVersion }}%s
%s%s
├── OS        : {{ .OS }}
├── NumCPU    : {{ .NumCPU }}
├── Compiler  : {{ .Compiler }}
└── Date      : {{ Now "Monday, 2 Jan 2006" }}%s
`

const shortVersionBanner string = ` _____     _ 
/ _  / ___| |
\// / / _ \ |
 / //\  __/ |__
/____/\___|___/  v{{ .ZelVersion }}
`

var CmdVersion = &commands.Command{
	UsageLine: "version",
	Short:     "Prints the current Zel version",
	Long: `
Prints the current Bee, Beego and Go version alongside the platform information.
`,
	Run: versionCmd,
}

var outputFormat string

const version = config.Version

func init() {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.StringVar(&outputFormat, "o", "", "Set the output format. Either json or yaml.")
	CmdVersion.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdVersion)
}

func versionCmd(cmd *commands.Command, args []string) int {
	cmd.Flag.Parse(args)
	stdout := cmd.Out()

	if outputFormat != "" {
		runtimeInfo := RuntimeInfo{
			OS:         runtime.GOOS,
			NumCPU:     runtime.NumCPU(),
			Compiler:   runtime.Compiler,
			ZelVersion: version,
		}
		switch outputFormat {
		case "json":
			{
				b, err := json.MarshalIndent(runtimeInfo, "", "    ")
				if err != nil {
					logger.Log.Error(err.Error())
				}
				fmt.Println(string(b))
				return 0
			}
		case "yaml":
			{
				b, err := yaml.Marshal(&runtimeInfo)
				if err != nil {
					logger.Log.Error(err.Error())
				}
				fmt.Println(string(b))
				return 0
			}
		}
	}

	coloredBanner := fmt.Sprintf(verboseVersionBanner, "\x1b[35m", "\x1b[1m",
		"\x1b[0m", "\x1b[32m", "\x1b[1m", "\x1b[0m")
	InitBanner(stdout, bytes.NewBufferString(coloredBanner))
	return 0
}

func ShowShortVersionBanner() {
	output := colors.NewColorWriter(os.Stdout)
	InitBanner(output, bytes.NewBufferString(colors.MagentaBold(shortVersionBanner)))
}
