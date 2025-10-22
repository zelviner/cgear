package cmd

import (
	"github.com/zelviner/cgear/cmd/commands"
	_ "github.com/zelviner/cgear/cmd/commands/build"
	_ "github.com/zelviner/cgear/cmd/commands/count"
	_ "github.com/zelviner/cgear/cmd/commands/env"
	_ "github.com/zelviner/cgear/cmd/commands/generate"
	_ "github.com/zelviner/cgear/cmd/commands/install"
	_ "github.com/zelviner/cgear/cmd/commands/new"
	_ "github.com/zelviner/cgear/cmd/commands/pack"
	_ "github.com/zelviner/cgear/cmd/commands/run"
	_ "github.com/zelviner/cgear/cmd/commands/test"
	_ "github.com/zelviner/cgear/cmd/commands/version"
	"github.com/zelviner/cgear/utils"
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

var usageTemplate = `Cgear is a Fast tool for managing your C++ Project.

{{"USAGE" | headline}}
    {{"cgear command [arguments]" | bold}}

{{"AVAILABLE COMMANDS" | headline}}
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-11s" | bold}} {{.Short}}{{end}}{{end}}

Use {{"cgear help [command]" | bold}} for more information about a command.

{{"ADDITIONAL HELP TOPICS" | headline}}
{{range .}}{{if not .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use {{"cgear help [topic]" | bold}} for more information about that topic.
`

var helpTemplate = `{{"USAGE" | headline}}
  {{.UsageLine | printf "cgear %s" | bold}}
{{if .Options}}{{endline}}{{"OPTIONS" | headline}}{{range $k,$v := .Options}}
  {{$k | printf "-%s" | bold}}
      {{$v}}
  {{end}}{{end}}
{{"DESCRIPTION" | headline}}
  {{tmpltostr .Long . | trim}}
`

var ErrorTemplate = `cgear: %s.
Use {{"cgear help" | bold}} for more information.
`

// cgear tool 使用说明
func Usage() {
	utils.Tmpl(usageTemplate, commands.AvailableCommands)
}

// cgear tool 帮助信息
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
