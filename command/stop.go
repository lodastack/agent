package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"syscall"

	"github.com/lodastack/agent/config"

	"github.com/oiooj/cli"
)

var CmdStop = cli.Command{
	Name:        "stop",
	Usage:       "关闭客户端",
	Description: "关闭Agent客户端",
	Action:      runStop,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "m",
			Value: "normal",
			Usage: "stop agent mode",
		},
	},
}

func runStop(c *cli.Context) {
	stopAgent()
	if c.String("m") == "clean" {
		cleanData()
	}
}

func stopAgent() {
	if runtime.GOOS != "linux" {
		fmt.Printf("Agent don't support this arch: %s\n", runtime.GOOS)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(config.PID)
	if err != nil {
		fmt.Printf("cannot read pid file: %s ", config.PID)
		return
	}
	pid, _ := strconv.Atoi(string(data))
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("read pid process error ")
		return
	}
	err = process.Signal(syscall.SIGINT)
	if err != nil {
		fmt.Printf("send signal to process error: %s ", err)
		return
	}
	err = os.Remove(config.PID)
	if err == os.ErrNotExist {
		err = nil
	}
	if err != nil {
		fmt.Printf("send signal to process successfully and remove pid error: %s ", err)
	}
	fmt.Printf("send signal to process successfully ")
}

func cleanData() {
	// TODO:stop don't know the data dir, data dir is a var, need to use signal to pass remove cmd.
	if err := os.RemoveAll("/usr/local/agent-plugins"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("clean agent data finish")
}
