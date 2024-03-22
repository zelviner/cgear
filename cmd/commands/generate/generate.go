package generate

import (
	"os"
	"zel/cmd/commands"
	"zel/cmd/commands/version"
	"zel/generate"
	"zel/logger"
	"zel/utils"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var CmdGenerate = &commands.Command{
	UsageLine: "generate [command]",
	Short:     "Source code generator",
	Long: `▶ {{"To scaffold out your entire application:"|bold}}

     $ bee generate scaffold [scaffoldname] [-fields="title:string,body:text"] [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"]

  ▶ {{"To generate a Model based on fields:"|bold}}

     $ bee generate model [modelname] [-fields="name:type"]

  ▶ {{"To generate a controller:"|bold}}

     $ bee generate controller [controllerfile]

  ▶ {{"To generate a CRUD view:"|bold}}

     $ bee generate view [viewpath]

  ▶ {{"To generate a migration file for making database schema updates:"|bold}}

     $ bee generate migration [migrationfile] [-fields="name:type"]

  ▶ {{"To generate swagger doc file:"|bold}}

     $ bee generate docs

    ▶ {{"To generate swagger doc file:"|bold}}

     $ bee generate routers [-ctrlDir=/path/to/controller/directory] [-routersFile=/path/to/routers/file.go] [-routersPkg=myPackage]

  ▶ {{"To generate a test case:"|bold}}

     $ bee generate test [routerfile]

  ▶ {{"To generate appcode based on an existing database:"|bold}}

     $ bee generate appcode [-tables=""] [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"] [-level=3]
	`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    GenerateCode,
}

func init() {
	CmdGenerate.Flag.Var(&generate.Include, "test", "Generate C++ project header folder.")

	commands.AvailableCommands = append(commands.AvailableCommands, CmdGenerate)
}

func GenerateCode(cmd *commands.Command, args []string) int {
	currPath, _ := os.Getwd()
	if len(args) < 1 {
		logger.Log.Fatal("Command is missing")
	}

	gcmd := args[0]
	switch gcmd {
	case "include":
		include(cmd, args, currPath)

	default:
		logger.Log.Fatal("Command is missing")
	}

	logger.Log.Successf("%s successfully generated!", cases.Title(language.English).String(gcmd))

	return 0
}

func include(cmd *commands.Command, args []string, currPath string) {
	utils.IsZelProject(currPath)
}
