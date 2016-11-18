package sysinfo

import (
	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

func KernelMetrics() (L []*common.Metric) {
	maxFiles, err := nux.KernelMaxFiles()
	if err != nil {
		log.Error("failed collect kernel metrics:", err)
		return
	}

	L = append(L, toMetric("kernel.maxfiles", maxFiles, nil))

	maxProc, err := nux.KernelMaxProc()
	if err != nil {
		log.Error(err)
		return
	}

	L = append(L, toMetric("kernel.maxproc", maxProc, nil))

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Error("failed to call KernelAllocateFiles:", err)
		return
	}

	L = append(L, toMetric("kernel.files.allocated", allocateFiles, nil))
	L = append(L, toMetric("kernel.files.left", maxFiles-allocateFiles, nil))
	return
}
