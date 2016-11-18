package sysinfo

import (
	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

func SocketStatSummaryMetrics() (L []*common.Metric) {
	ssMap, err := nux.SocketStatSummary()
	if err != nil {
		log.Error("failed to collect SocketStatSummaryMetrics:", err)
		return
	}

	for k, v := range ssMap {
		L = append(L, toMetric("net."+k, v, nil))
	}

	return
}
