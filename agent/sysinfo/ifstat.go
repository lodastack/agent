package sysinfo

import (
	"time"

	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

const MILLION_BIT = 1000000
const BITS_PER_BYTE = 8

type CumIfStat struct {
	inBytes  int64
	outBytes int64
	inDrop   int64
	outDrop  int64
	speed    int64
}

var (
	historyIfStat map[string]CumIfStat
	lastTime      time.Time
)

func NetMetrics() (ret []*common.Metric) {
	netIfs, err := nux.NetIfs(common.Conf.IfacePrefix)
	if err != nil {
		log.Error("collect net metric accurs error:", err)
		return
	}
	now := time.Now()
	newIfStat := make(map[string]CumIfStat)
	for _, netIf := range netIfs {
		newIfStat[netIf.Iface] = CumIfStat{netIf.InBytes, netIf.OutBytes, netIf.InDropped, netIf.OutDropped, netIf.Speed}
	}
	interval := now.Unix() - lastTime.Unix()
	lastTime = now

	if historyIfStat != nil {
		for iface, stat := range newIfStat {
			tags := map[string]string{"interface": iface}
			oldStat := historyIfStat[iface]
			netIn := common.SetPrecision(float64(stat.inBytes-oldStat.inBytes)*BITS_PER_BYTE/float64(interval), 2)
			ret = append(ret, toMetric("net.in", netIn, tags))

			netOut := common.SetPrecision(float64(stat.outBytes-oldStat.outBytes)*BITS_PER_BYTE/float64(interval), 2)
			ret = append(ret, toMetric("net.out", netOut, tags))

			v := common.SetPrecision(float64(stat.inDrop-oldStat.inDrop)/float64(interval), 2)
			ret = append(ret, toMetric("net.in.dropped", v, tags))

			v = common.SetPrecision(float64(stat.outDrop-oldStat.outDrop)/float64(interval), 2)
			ret = append(ret, toMetric("net.out.dropped", v, tags))

			if stat.speed != 0 {
				v = common.SetPrecision(float64(netIn*100/float64(stat.speed*MILLION_BIT)), 2)
				ret = append(ret, toMetric("net.in.percent", v, tags))

				v = common.SetPrecision(float64(netOut*100/float64(stat.speed*MILLION_BIT)), 2)
				ret = append(ret, toMetric("net.out.percent", v, tags))
			}

			ret = append(ret, toMetric("net.speed", stat.speed, tags))
		}

	}
	historyIfStat = newIfStat
	return
}
