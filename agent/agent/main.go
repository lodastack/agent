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
}

// New returns an Agent struct based off the given Config.
func New(c *config.Config) (*Agent, error) {
	a := &Agent{
		Config: &c.Agent,
	}

	var err error
	a.Output, err = outputs.New(&c.Output)
	return a, err
}

// Start starts the agent collects data.
func (a *Agent) Start() {
	common.InitCollectConfig(a.Config)

	go a.Output.Start()
	go scheduler.Start()
	go httpd.Start(a.Config.Listen)
	go a.Report()
}
