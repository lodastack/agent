package command

import (
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"

	"github.com/lodastack/agent/agent/agent"
	"github.com/lodastack/agent/config"
	"github.com/lodastack/log"

	"github.com/oiooj/cli"
)

var CmdDebug = cli.Command{
	Name:        "debug",
	Usage:       "调试模式",
	Description: "调试模式",
	Action:      runDebug,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "f",
			Value: "/etc/monitor-agent.conf",
			Usage: "配置文件路径，默认位置：/etc/monitor-agent.conf",
		},
		cli.StringFlag{
			Name:  "cpuprofile",
			Value: "/tmp/cpu.pprof",
			Usage: "CPU pprof file",
		},
		cli.StringFlag{
			Name:  "memprofile",
			Value: "/tmp/mem.pprof",
			Usage: "Memory pprof file",
		},
	},
}

func runDebug(c *cli.Context) {
	// parse config file
	err := config.ParseConfig(c.String("f"))
	if err != nil {
		log.Fatalf("Parse Config File Error : " + err.Error())
	}
	// init dlog setting
	initLog()
	// save pid to file
	ioutil.WriteFile(config.PID, []byte(strconv.Itoa(os.Getpid())), 0744)
	go Notify()

	//start agent module
	a, err := agent.New(config.C)
	if err != nil {
		log.Fatalf("New agent Error : " + err.Error())
	}
	a.Start()

	// start pprof
	startProfile(c.String("cpuprofile"), c.String("memprofile"))
	select {}
}

// prof stores the file locations of active profiles.
var prof struct {
	cpu *os.File
	mem *os.File
}

// startProfile initializes the CPU and memory profile, if specified.
func startProfile(cpuprofile, memprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Errorf("failed to create CPU profile file at %s: %s", cpuprofile, err.Error())
		}
		log.Printf("writing CPU profile to: %s\n", cpuprofile)
		prof.cpu = f
		pprof.StartCPUProfile(prof.cpu)
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Errorf("failed to create memory profile file at %s: %s", cpuprofile, err.Error())
		}
		log.Printf("writing memory profile to: %s\n", memprofile)
		prof.mem = f
		runtime.MemProfileRate = 4096
	}
}

// stopProfile closes the CPU and memory profiles if they are running.
func stopProfile() {
	if prof.cpu != nil {
		pprof.StopCPUProfile()
		prof.cpu.Close()
		log.Printf("CPU profiling stopped")
	}
	if prof.mem != nil {
		pprof.Lookup("heap").WriteTo(prof.mem, 0)
		prof.mem.Close()
		log.Printf("memory profiling stopped")
	}
}
