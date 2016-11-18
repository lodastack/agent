package sysinfo

import (
	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

var USES = map[string]struct{}{
	"PruneCalled":        {},
	"LockDroppedIcmps":   {},
	"ArpFilter":          {},
	"TW":                 {},
	"DelayedACKLocked":   {},
	"ListenOverflows":    {},
	"ListenDrops":        {},
	"TCPPrequeueDropped": {},
	"TCPTSReorder":       {},
	"TCPDSACKUndo":       {},
	"TCPLoss":            {},
	"TCPLostRetransmit":  {},
	"TCPLossFailures":    {},
	"TCPFastRetrans":     {},
	"TCPTimeouts":        {},
	"TCPSchedulerFailed": {},
	"TCPAbortOnMemory":   {},
	"TCPAbortOnTimeout":  {},
	"TCPAbortFailed":     {},
	"TCPMemoryPressures": {},
	"TCPSpuriousRTOs":    {},
	"TCPBacklogDrop":     {},
	"TCPMinTTLDrop":      {},
}

func NetstatMetrics() (L []*common.Metric) {
	tcpExts, err := nux.Netstat("TcpExt")

	if err != nil {
		log.Error("failed to collect NetstatMetrics:", err)
		return
	}

	cnt := len(tcpExts)
	if cnt == 0 {
		return
	}

	for key, val := range tcpExts {
		if _, ok := USES[key]; !ok {
			continue
		}
		L = append(L, toMetric("TcpExt."+key, val, nil))
	}

	return
}
