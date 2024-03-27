package main

import (
	"flag"
	"log"
	"os"

	"github.com/ZEL-30/zel/cmd"
	"github.com/ZEL-30/zel/cmd/commands"
	"github.com/ZEL-30/zel/config"
	"github.com/ZEL-30/zel/utils"
)

func main() {

	flag.Usage = cmd.Usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()

	if len(args) < 1 {
		cmd.Usage()
		os.Exit(2)
		return
	}

	if args[0] == "help" {
		cmd.Help(args[1:])
		return
	}

	for _, c := range commands.AvailableCommands {
		if c.Name() == args[0] && c.Run != nil {
			c.Flag.Usage = func() { c.Usage() }
			if c.CustomFlags {
				args = args[1:]
			} else {
				c.Flag.Parse(args[1:])
				args = c.Flag.Args()
			}

			if c.PreRun != nil {
				c.PreRun(c, args)
			}

			config.LaodConfig()
			os.Exit(c.Run(c, args))
			return
		}

	}

	utils.PrintErrorAndExit("Unknow subcommand", cmd.ErrorTemplate)
}
