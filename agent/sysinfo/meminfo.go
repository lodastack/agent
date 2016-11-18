package sysinfo

import (
	"math"

	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

func MemMetrics() []*common.Metric {
	m, err := nux.MemInfo()
	if err != nil {
		log.Error("failed to collect Metrics:", err)
		return nil
	}

	memFree := m.MemFree + m.Buffers + m.Cached
	memUsed := m.MemTotal - memFree

	pmemUsed := 0.0
	if m.MemTotal != 0 {
		pmemUsed = math.Ceil(float64(memUsed) * 100.0 / float64(m.MemTotal))
	}

	return []*common.Metric{
		toMetric("mem.total", m.MemTotal, nil),
		toMetric("mem.used", memUsed, nil),
		toMetric("mem.free", memFree, nil),
		toMetric("mem.used.percent", pmemUsed, nil),
		toMetric("mem.buffers", m.Buffers, nil),
		toMetric("mem.cached", m.Cached, nil),
	}

}
