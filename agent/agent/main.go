package agent

import (
	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/httpd"
	"github.com/lodastack/agent/agent/outputs"
	_ "github.com/lodastack/agent/agent/outputs/all"
	"github.com/lodastack/agent/agent/scheduler"
	"github.com/lodastack/agent/config"
)

// Agent runs collects data based on the given config.
type Agent struct {
	Config *common.AgentConfig
	Output *outputs.Output
	Httpd  *httpd.Service
}

// New returns an Agent struct based off the given Config.
func New(c *config.Config) (*Agent, error) {
	a := &Agent{
		Config: &c.Agent,
		Httpd:  httpd.NewService(c.Agent.Listen),
	}

	var err error
	a.Output, err = outputs.New(&c.Output)
	return a, err
}

// Start starts the agent collects data.
func (a *Agent) Start() error {
	common.InitCollectConfig(a.Config)

	go a.Output.Start()
	go scheduler.Start()
	if err := a.Httpd.Start(); err != nil {
		return err
	}
	go a.Report()
	return nil
}
