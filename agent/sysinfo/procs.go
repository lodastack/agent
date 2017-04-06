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
	CPUTotal     map[int]float64
	CPUProc      map[int]float64
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
	newCPUTotal := make(map[int]float64)
	newCPUProc := make(map[int]float64)
	interval := float64(time.Now().Unix() - lastProcTime.Unix())
	lastProcTime = time.Now()
	for _, p := range ps {
		newRBytes[p.Pid] = p.RBytes
		newWBytes[p.Pid] = p.WBytes
		newCPUTotal[p.Pid] = p.TotalCpu
		newCPUProc[p.Pid] = p.Cpu
	}
	for _, proc := range procs {
		var cnt int
		var fdNum int
		var memory uint64
		var cpu, procCpuTime, totalCpuTime float64
		var ioWrite, ioRead uint64
		for _, p := range ps {
			if proc.BinaryPath == p.Exe {
				cnt++
				memory += p.Mem

				if totalCpuOld, _ := CPUTotal[p.Pid]; totalCpuOld <= p.TotalCpu {
					totalCpuTime = p.TotalCpu - totalCpuOld
				}
				if cpuOld, _ := CPUProc[p.Pid]; cpuOld <= p.Cpu {
					procCpuTime = p.Cpu - cpuOld
				}
				if totalCpuTime > 0 {
					cpu += procCpuTime / totalCpuTime
				}

				fdNum += p.FdCount
				if rOld, _ := rBytes[p.Pid]; rOld <= p.RBytes {
					ioRead += p.RBytes - rOld
				}
				if wOld, _ := wBytes[p.Pid]; wOld <= p.WBytes {
					ioWrite += p.WBytes - wOld
				}
			}
		}
		m[proc.Namespace] = append(m[proc.Namespace],
			toMetric(fmt.Sprintf("%s.procnum", proc.Name), cnt, nil),
			toMetric(fmt.Sprintf("%s.fdnum", proc.Name), fdNum, nil),
			// unit:Byte
			toMetric(fmt.Sprintf("%s.mem", proc.Name), memory*1024, nil),
			toMetric(fmt.Sprintf("%s.cpu", proc.Name), common.SetPrecision(cpu*100, 2), nil))

		if rBytes != nil {
			m[proc.Namespace] = append(m[proc.Namespace],
				// unit:Byte
				toMetric(fmt.Sprintf("%s.io.read", proc.Name), math.Ceil(float64(ioRead)/(interval)), nil),
				toMetric(fmt.Sprintf("%s.io.write", proc.Name), math.Ceil(float64(ioWrite)/(interval)), nil))
		}
	}
	rBytes = newRBytes
	wBytes = newWBytes
	CPUTotal = newCPUTotal
	CPUProc = newCPUProc
	for k, v := range m {
		outputs.SendMetrics(common.TYPE_PROC, k, v)
	}
	return
}

func (self ProcCollector) Description() string {
	return "ProcessCollector"
}
