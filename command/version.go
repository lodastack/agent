package command

import (
	"fmt"
	"os"
	"runtime"

	"github.com/lodastack/agent/config"

	"github.com/oiooj/cli"
)

var CmdVersion = cli.Command{
	Name:        "version",
	Usage:       "show version",
	Description: "show version",
	Action:      runVersion,
}

func runVersion(c *cli.Context) {
	// Print version info.
	fmt.Fprintf(os.Stdout, "Monitor Agent v%s (git: %s %s) build: %s %s\n", config.Version, config.Branch, config.Commit, config.BuildTime, runtime.Version())
}
