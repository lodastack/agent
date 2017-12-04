package sysinfo

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

var (
	diskStatsMap = make(map[string][2]*nux.DiskStats)
	dsLock       = new(sync.RWMutex)
)

func UpdateDiskStats() error {
	dsList, err := nux.ListDiskStats()
	if err != nil {
		log.Error("failed to UpdateDiskStats:", err)
		return err
	}

	dsLock.Lock()
	defer dsLock.Unlock()
	for i := 0; i < len(dsList); i++ {
		device := dsList[i].Device
		diskStatsMap[device] = [2]*nux.DiskStats{dsList[i], diskStatsMap[device][0]}
	}
	return nil
}

func IOReadRequests(arr [2]*nux.DiskStats) uint64 {
	return arr[0].ReadRequests - arr[1].ReadRequests
}

func IOReadMerged(arr [2]*nux.DiskStats) uint64 {
	return arr[0].ReadMerged - arr[1].ReadMerged
}

func IOReadSectors(arr [2]*nux.DiskStats) uint64 {
	return arr[0].ReadSectors - arr[1].ReadSectors
}

func IOMsecRead(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecRead - arr[1].MsecRead
}

func IOWriteRequests(arr [2]*nux.DiskStats) uint64 {
	return arr[0].WriteRequests - arr[1].WriteRequests
}

func IOWriteMerged(arr [2]*nux.DiskStats) uint64 {
	return arr[0].WriteMerged - arr[1].WriteMerged
}

func IOWriteSectors(arr [2]*nux.DiskStats) uint64 {
	return arr[0].WriteSectors - arr[1].WriteSectors
}

func IOMsecWrite(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecWrite - arr[1].MsecWrite
}

func IOMsecTotal(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecTotal - arr[1].MsecTotal
}

func IOMsecWeightedTotal(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecWeightedTotal - arr[1].MsecWeightedTotal
}

func TS(arr [2]*nux.DiskStats) uint64 {
	return uint64(arr[0].TS.Sub(arr[1].TS).Nanoseconds() / 1000000)
}

func IODelta(device string, f func([2]*nux.DiskStats) uint64) uint64 {
	val, ok := diskStatsMap[device]
	if !ok {
		return 0
	}

	if val[1] == nil {
		return 0
	}
	return f(val)
}

func DiskIOMetrics() (L []*common.Metric) {
	dsList, err := nux.ListDiskStats()
	if err != nil {
		log.Error("error when collect disk io:", err)
		return
	}

	for _, ds := range dsList {
		if !ShouldHandleDevice(ds.Device) {
			continue
		}

		tags := map[string]string{"device": ds.Device}

		L = append(L, toMetric("disk.io.read_requests", ds.ReadRequests, tags))
		L = append(L, toMetric("disk.io.read_merged", ds.ReadMerged, tags))
		L = append(L, toMetric("disk.io.read_sectors", ds.ReadSectors, tags))
		L = append(L, toMetric("disk.io.msec_read", ds.MsecRead, tags))
		L = append(L, toMetric("disk.io.write_requests", ds.WriteRequests, tags))
		L = append(L, toMetric("disk.io.write_merged", ds.WriteMerged, tags))
		L = append(L, toMetric("disk.io.write_sectors", ds.WriteSectors, tags))
		L = append(L, toMetric("disk.io.msec_write", ds.MsecWrite, tags))
		L = append(L, toMetric("disk.io.ios_in_progress", ds.IosInProgress, tags))
		L = append(L, toMetric("disk.io.msec_total", ds.MsecTotal, tags))
		L = append(L, toMetric("disk.io.msec_weighted_total", ds.MsecWeightedTotal, tags))
	}
	return
}

func IOStatsMetrics() (L []*common.Metric) {
	dsLock.RLock()
	defer dsLock.RUnlock()

	for device := range diskStatsMap {
		if !ShouldHandleDevice(device) {
			continue
		}

		tags := map[string]string{"device": device}
		rio := IODelta(device, IOReadRequests)
		wio := IODelta(device, IOWriteRequests)
		ruse := IODelta(device, IOMsecRead)
		wuse := IODelta(device, IOMsecWrite)
		use := IODelta(device, IOMsecTotal)
		n_io := rio + wio
		await := 0.0
		if n_io != 0 {
			await = common.SetPrecision(float64(ruse+wuse)/float64(n_io), 2)
		}

		duration := IODelta(device, TS)
		L = append(L, toMetric("disk.io.read_requests", common.SetPrecision(float64(rio)/float64(duration/1000), 2), tags))
		L = append(L, toMetric("disk.io.write_requests", common.SetPrecision(float64(wio)/float64(duration/1000), 2), tags))
		L = append(L, toMetric("disk.io.await", await, tags))
		tmp := common.SetPrecision(float64(use)*100.0/float64(duration), 2)
		if tmp > 100.0 {
			tmp = 100.0
		}
		if tmp >= 0 {
			L = append(L, toMetric("disk.io.util", tmp, tags))
		}
	}

	return
}

func IOStatsForPage() (L [][]string) {
	dsLock.RLock()
	defer dsLock.RUnlock()

	for device := range diskStatsMap {
		if !ShouldHandleDevice(device) {
			continue
		}

		rio := IODelta(device, IOReadRequests)
		wio := IODelta(device, IOWriteRequests)

		delta_rsec := IODelta(device, IOReadSectors)
		delta_wsec := IODelta(device, IOWriteSectors)

		ruse := IODelta(device, IOMsecRead)
		wuse := IODelta(device, IOMsecWrite)
		use := IODelta(device, IOMsecTotal)
		n_io := rio + wio
		avgrq_sz := 0.0
		await := 0.0
		svctm := 0.0
		if n_io != 0 {
			avgrq_sz = float64(delta_rsec+delta_wsec) / float64(n_io)
			await = float64(ruse+wuse) / float64(n_io)
			svctm = float64(use) / float64(n_io)
		}

		item := []string{
			device,
			fmt.Sprintf("%d", IODelta(device, IOReadMerged)),
			fmt.Sprintf("%d", IODelta(device, IOWriteMerged)),
			fmt.Sprintf("%d", rio),
			fmt.Sprintf("%d", wio),
			fmt.Sprintf("%.2f", float64(delta_rsec)/2.0),
			fmt.Sprintf("%.2f", float64(delta_wsec)/2.0),
			fmt.Sprintf("%.2f", avgrq_sz),                                             // avgrq-sz: delta(rsect+wsect)/delta(rio+wio)
			fmt.Sprintf("%.2f", float64(IODelta(device, IOMsecWeightedTotal))/1000.0), // avgqu-sz: delta(aveq)/s/1000
			fmt.Sprintf("%.2f", await),                                                // await: delta(ruse+wuse)/delta(rio+wio)
			fmt.Sprintf("%.2f", svctm),                                                // svctm: delta(use)/delta(rio+wio)
			fmt.Sprintf("%.2f%%", float64(use)/10.0),                                  // %util: delta(use)/s/1000 * 100%
		}
		L = append(L, item)
	}

	return
}

func ShouldHandleDevice(device string) bool {
	normal := len(device) == 3 && (strings.HasPrefix(device, "sd") || strings.HasPrefix(device, "vd"))
	aws := len(device) == 4 && strings.HasPrefix(device, "xvd")
	fusion := len(device) == 4 && (strings.HasPrefix(device, "hio") || strings.HasPrefix(device, "fio"))
	return normal || aws || fusion
}
