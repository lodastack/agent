package sysinfo

import (
	"strconv"
	"sync"

	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

const (
	historyCount int = 2
)

var (
	procStatHistory [historyCount]*nux.ProcStat
	psLock          = new(sync.RWMutex)
)

func UpdateCpuStat() error {
	ps, err := nux.CurrentProcStat()
	if err != nil {
		log.Error("failed to UpdateCpuStat:", err)
		return err
	}

	psLock.Lock()
	defer psLock.Unlock()
	for i := historyCount - 1; i > 0; i-- {
		procStatHistory[i] = procStatHistory[i-1]
	}

	procStatHistory[0] = ps
	return nil
}

func deltaTotal() uint64 {
	if procStatHistory[1] == nil {
		return 0
	}
	return procStatHistory[0].Cpu.Total - procStatHistory[1].Cpu.Total
}

func CpuIdle() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return common.SetPrecision(float64(procStatHistory[0].Cpu.Idle-procStatHistory[1].Cpu.Idle)*invQuotient, 2)
}

func CpuIdles() (res []float64) {
	psLock.RLock()
	defer psLock.RUnlock()
	if procStatHistory[1] == nil {
		return
	}
	if len(procStatHistory[0].Cpus) != len(procStatHistory[1].Cpus) {
		return
	}
	for i, c := range procStatHistory[0].Cpus {
		dt := c.Total - procStatHistory[1].Cpus[i].Total
		if dt == 0 {
			return
		}
		invQuotient := 100.00 / float64(dt)
		res = append(res, common.SetPrecision(float64(c.Idle-procStatHistory[1].Cpus[i].Idle)*invQuotient, 2))
	}
	return
}

func CpuUser() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.User-procStatHistory[1].Cpu.User) * invQuotient
}

func CpuNice() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Nice-procStatHistory[1].Cpu.Nice) * invQuotient
}

func CpuSystem() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.System-procStatHistory[1].Cpu.System) * invQuotient
}

func CpuIowait() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Iowait-procStatHistory[1].Cpu.Iowait) * invQuotient
}

func CpuIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Irq-procStatHistory[1].Cpu.Irq) * invQuotient
}

func CpuSoftIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.SoftIrq-procStatHistory[1].Cpu.SoftIrq) * invQuotient
}

func CpuSteal() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Steal-procStatHistory[1].Cpu.Steal) * invQuotient
}

func CpuGuest() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Guest-procStatHistory[1].Cpu.Guest) * invQuotient
}

func CurrentCpuSwitches() uint64 {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[0].Ctxt
}

func CpuPrepared() bool {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[1] != nil
}

func CpuMetrics() []*common.Metric {
	if !CpuPrepared() {
		return []*common.Metric{}
	}

	cpuIdleVal := CpuIdle()
	idle := toMetric("cpu.idle", cpuIdleVal, nil)
	res := []*common.Metric{idle}

	idles := CpuIdles()
	for i, v := range idles {
		tags := map[string]string{"core": strconv.Itoa(i)}
		res = append(res, toMetric("cpu.idle.core", v, tags))
	}

	load, err := nux.LoadAvg()
	if err != nil {
		log.Error("failed to collect LoadAvgMetrics:", err)
	} else {
		res = append(res, toMetric("cpu.loadavg.1", load.Avg1min, nil))
		res = append(res, toMetric("cpu.loadavg.5", load.Avg5min, nil))
		res = append(res, toMetric("cpu.loadavg.15", load.Avg15min, nil))
	}
	return res
}
