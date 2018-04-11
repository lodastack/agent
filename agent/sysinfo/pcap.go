package sysinfo

import "time"

const (
	snapshotLen   int32         = 1024
	promiscuous   bool          = false
	pcaptimeout   time.Duration = -1 * time.Second
	runDuration   int           = 50
	maxPacketSize               = 30
)
