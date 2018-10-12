package sysinfo

import (
	"github.com/lodastack/agent/agent/common"
)

// FsKernelMetrics collects file system metrices
func FsKernelMetrics() (L []*common.Metric) {
	return nil
}

// PsMetrics exec `ps` to get all process states
func PsMetrics() (L []*common.Metric) {
	return nil
}

// WtmpMetrics collect users login history
func WtmpMetrics() (L []*common.Metric) {
	return nil
}

// BtmpMetrics collect users failed login history
func BtmpMetrics() (L []*common.Metric) {
	return nil
}
