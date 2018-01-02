package sysinfo

import (
	"github.com/lodastack/agent/agent/common"
)

// AgentMetrics report agent alive metric
func AgentMetrics() []*common.Metric {
	return []*common.Metric{toMetric("agent.alive", 1, nil)}
}
