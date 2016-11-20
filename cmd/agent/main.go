package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/lodastack/agent/command"
	"github.com/lodastack/agent/config"

	"github.com/oiooj/cli"
)

func init() {
	if runtime.GOOS != "linux" {
		fmt.Printf("Agent don't support this arch: %s\n", runtime.GOOS)
		os.Exit(1)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = config.AppName
	app.Usage = config.Usage
	app.Version = config.Version
	app.Author = config.Author
	app.Email = config.Email

	app.Commands = []cli.Command{
		command.CmdStart,
		command.CmdStop,
		command.CmdDebug,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
