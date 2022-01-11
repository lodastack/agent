package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/lodastack/agent/command"
	"github.com/lodastack/agent/config"

	"github.com/oiooj/cli"
)

// These variables are populated via the Go linker.
var (
	version   string
	commit    string
	branch    string
	buildTime string
)

func init() {
	// If commit, branch, or build time are not set, make that clear.
	config.Version = version
	if version == "" {
		config.Version = "unknown"
	}
	config.Commit = commit
	if commit == "" {
		config.Commit = "unknown"
	}
	config.Branch = branch
	if branch == "" {
		config.Branch = "unknown"
	}
	config.BuildTime = buildTime
	if buildTime == "" {
		config.BuildTime = "unknown"
	}
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	if runtime.GOOS == "windows" {
		command.WindowsStart()
		return
	}
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = config.AppName
	app.Usage = config.Usage
	app.Version = config.Version
	app.Author = config.Author
	app.Email = config.Email

	app.Commands = []cli.Command{
		command.CmdCloudStart,
		command.CmdStart,
		command.CmdStop,
		command.CmdVersion,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
