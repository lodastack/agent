package sysinfo

import (
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
		pmemUsed = common.SetPrecision(float64(memUsed)*100.0/float64(m.MemTotal), 2)
	}

	pswapUsed := 0.0
	if m.SwapTotal != 0 {
		pswapUsed = common.SetPrecision(float64(m.SwapUsed)*100.0/float64(m.SwapTotal), 2)
	}

	return []*common.Metric{
		toMetric("mem.total", m.MemTotal, nil),
		toMetric("mem.used", memUsed, nil),
		toMetric("mem.free", memFree, nil),
		toMetric("mem.used.percent", pmemUsed, nil),
		toMetric("mem.buffers", m.Buffers, nil),
		toMetric("mem.cached", m.Cached, nil),
		toMetric("mem.swap.used.percent", pswapUsed, nil),
	}

}
