package command

import (
	"fmt"
	"io/ioutil"
	"os"
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
}

func runStop(c *cli.Context) {
	stopAgent()
}

func stopAgent() {
	data, err := ioutil.ReadFile(config.PID)
	if err != nil {
		fmt.Printf("cannot read pid file: %s\n", config.PID)
		return
	}
	pid, _ := strconv.Atoi(string(data))
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("read pid process error")
		return
	}
	err = process.Signal(syscall.SIGINT)
	if err != nil {
		fmt.Println("send signal to process error")
		return
	}
	err = os.Remove(config.PID)
	if err == os.ErrNotExist {
		err = nil
	}
	if err != nil {
		fmt.Println("send signal to process successfully and remove pid error: ", err)
	}
	fmt.Println("send signal to process successfully")
}
