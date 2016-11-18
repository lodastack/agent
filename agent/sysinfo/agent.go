package sysinfo

import (
	"github.com/lodastack/agent/agent/common"
)

func AgentMetrics() []*common.Metric {
	return []*common.Metric{toMetric("agent.alive", 1, nil)}
}
