package cmd

import (
	"github.com/ZEL-30/zel/cmd/commands"
	_ "github.com/ZEL-30/zel/cmd/commands/build"
	_ "github.com/ZEL-30/zel/cmd/commands/generate"
	_ "github.com/ZEL-30/zel/cmd/commands/kit"
	_ "github.com/ZEL-30/zel/cmd/commands/new"
	_ "github.com/ZEL-30/zel/cmd/commands/pack"
	_ "github.com/ZEL-30/zel/cmd/commands/run"
	_ "github.com/ZEL-30/zel/cmd/commands/test"
	_ "github.com/ZEL-30/zel/cmd/commands/version"
	"github.com/ZEL-30/zel/utils"
)

func IfGenerateDocs(name string, args []string) bool {
	if name != "generate" {
		return false
	}

	for _, arg := range args {
		if arg == "docs" {
			return true
		}
	}
	return false
}

var usageTemplate = `Zel is a Fast tool for managing your C++ Project.

{{"USAGE" | headline}}
    {{"zel command [arguments]" | bold}}

{{"AVAILABLE COMMANDS" | headline}}
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-11s" | bold}} {{.Short}}{{end}}{{end}}

Use {{"zel help [command]" | bold}} for more information about a command.

{{"ADDITIONAL HELP TOPICS" | headline}}
{{range .}}{{if not .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use {{"zel help [topic]" | bold}} for more information about that topic.
`

var helpTemplate = `{{"USAGE" | headline}}
  {{.UsageLine | printf "zel %s" | bold}}
{{if .Options}}{{endline}}{{"OPTIONS" | headline}}{{range $k,$v := .Options}}
  {{$k | printf "-%s" | bold}}
      {{$v}}
  {{end}}{{end}}
{{"DESCRIPTION" | headline}}
  {{tmpltostr .Long . | trim}}
`

var ErrorTemplate = `zel: %s.
Use {{"zel help" | bold}} for more information.
`

// zel tool 使用说明
func Usage() {
	utils.Tmpl(usageTemplate, commands.AvailableCommands)
}

// zel tool 帮助信息
func Help(args []string) {

	if len(args) == 0 {
		Usage()
		return
	}

	if len(args) != 1 {
		utils.PrintErrorAndExit("Too many arguments", ErrorTemplate)
	}

	arg := args[0]

	for _, cmd := range commands.AvailableCommands {
		if cmd.Name() == arg {
			utils.Tmpl(helpTemplate, cmd)
			return
		}
	}

	utils.PrintErrorAndExit("Unknown help topic", ErrorTemplate)
}
