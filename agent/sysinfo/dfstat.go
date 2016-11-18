package sysinfo

import (
	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

func DeviceMetrics() (L []*common.Metric) {
	mountPoints, err := nux.ListMountPoint()

	if err != nil {
		log.Error("failed to call ListMountPoint:", err)
		return
	}

	for idx := range mountPoints {
		var du *nux.DeviceUsage
		du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
		if err != nil {
			log.Error("failed to call BuildDeviceUsage:", err)
			continue
		}

		tags := map[string]string{"mount": du.FsFile}
		L = append(L, toMetric("disk.inodes.used.percent", du.InodesUsedPercent, tags))
		L = append(L, toMetric("disk.used.percent", du.BlocksUsedPercent, tags))
	}

	return
}
