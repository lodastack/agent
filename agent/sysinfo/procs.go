package sysinfo

import (
	"fmt"
	"math"
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/outputs"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

type ProcCollector struct{}

var (
	rBytes       map[int]uint64
	wBytes       map[int]uint64
	lastProcTime time.Time
)

func (self ProcCollector) Run() {
	procs := common.ReportProcs()
	cmdlines := map[string]string{}
	for _, proc := range procs {
		cmdlines[proc.BinaryPath] = proc.Name
	}
	ps, err := nux.Procs(cmdlines)
	if err != nil {
		log.Error("failed to collect ProcMetrics:", err)
		return
	}

	m := map[string][]*common.Metric{}

	newRBytes := make(map[int]uint64)
	newWBytes := make(map[int]uint64)
	interval := float64(time.Now().Unix() - lastProcTime.Unix())
	lastProcTime = time.Now()
	for _, p := range ps {
		newRBytes[p.Pid] = p.RBytes
		newWBytes[p.Pid] = p.WBytes
	}
	for _, proc := range procs {
		var cnt int
		var fdNum int
		var memory uint64
		var cpu float64
		var tcpEstablished int
		var ioWrite, ioRead uint64
		for _, p := range ps {
			if proc.BinaryPath == p.Exe {
				cnt++
				memory += p.Mem
				cpu += p.Cpu
				fdNum += p.FdCount
				tcpEstablished += p.TcpEstab
				if rOld := rBytes[p.Pid]; rOld <= p.RBytes {
					ioRead += p.RBytes - rOld
				}
				if wOld := wBytes[p.Pid]; wOld <= p.WBytes {
					ioWrite += p.WBytes - wOld
				}
			}
		}
		m[proc.Namespace] = append(m[proc.Namespace],
			toMetric(fmt.Sprintf("%s.%s.procnum", common.TYPE_PROC, proc.Name), cnt, nil),
			toMetric(fmt.Sprintf("%s.%s.fdnum", common.TYPE_PROC, proc.Name), fdNum, nil),
			// unit:Byte
			toMetric(fmt.Sprintf("%s.%s.mem", common.TYPE_PROC, proc.Name), memory*1024, nil),
			toMetric(fmt.Sprintf("%s.%s.cpu", common.TYPE_PROC, proc.Name), common.SetPrecision(cpu*100, 2), nil))

		if rBytes != nil {
			m[proc.Namespace] = append(m[proc.Namespace],
				// unit:Byte
				toMetric(fmt.Sprintf("%s.%s.io.read", common.TYPE_PROC, proc.Name), math.Ceil(float64(ioRead)/(interval)), nil),
				toMetric(fmt.Sprintf("%s.%s.io.write", common.TYPE_PROC, proc.Name), math.Ceil(float64(ioWrite)/(interval)), nil))
		}
	}
	rBytes = newRBytes
	wBytes = newWBytes
	for k, v := range m {
		outputs.SendMetrics(common.TYPE_PROC, k, v)
	}
	return
}

func (self ProcCollector) Description() string {
	return "ProcessCollector"
}
