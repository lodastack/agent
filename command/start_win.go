package command

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
	"github.com/lodastack/agent/agent/agent"
	"github.com/lodastack/agent/config"
	"github.com/lodastack/log"
)

const configFile = `C:\monitor-agent\agent.conf`

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	startAgent(configFile)
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func WindowsStart() {
	svcConfig := &service.Config{
		Name:        "MonitorAgent",
		DisplayName: "Monitor Agent",
		Description: "customer monitor service",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			s.Install()
			fmt.Println("服务安装成功")
			return
		}

		if os.Args[1] == "remove" {
			s.Uninstall()
			fmt.Println("服务卸载成功")
			return
		}
	}

	err = s.Run()
	if err != nil {
		panic(err)
	}
}

func startAgent(cf string) {
	//parse config file
	err := config.ParseConfig(cf)
	if err != nil {
		log.Fatalf("Parse Config File Error: %s", err.Error())
	}
	//init log setting
	initLog()
	//start agent module
	a, err := agent.New(config.C)
	if err != nil {
		log.Fatalf("New agent Error: %s", err.Error())
	}
	if err := a.Start(); err != nil {
		log.Fatalf("agent start failed: %s", err.Error())
	}
	// Print sweet Agent logo.
	PrintLogo()
	select {}
}
