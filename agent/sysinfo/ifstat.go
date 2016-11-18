package sysinfo

import (
	"time"

	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

type CumIfStat struct {
	inBytes  int64
	outBytes int64
	inDrop   int64
	outDrop  int64
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
		newIfStat[netIf.Iface] = CumIfStat{netIf.InBytes, netIf.OutBytes, netIf.InDropped, netIf.OutDropped}
	}
	interval := now.Unix() - lastTime.Unix()
	lastTime = now

	if historyIfStat != nil {
		for iface, stat := range newIfStat {
			tags := map[string]string{"interface": iface}
			oldStat := historyIfStat[iface]
			v := common.SetPrecision(float64(stat.inBytes-oldStat.inBytes)*8/float64(interval), 2)
			ret = append(ret, toMetric("net.in", v, tags))

			v = common.SetPrecision(float64(stat.outBytes-oldStat.outBytes)*8/float64(interval), 2)
			ret = append(ret, toMetric("net.out", v, tags))

			v = common.SetPrecision(float64(stat.inDrop-oldStat.inDrop)/float64(interval), 2)
			ret = append(ret, toMetric("net.in.dropped", v, tags))

			v = common.SetPrecision(float64(stat.outDrop-oldStat.outDrop)/float64(interval), 2)
			ret = append(ret, toMetric("net.out.dropped", v, tags))
		}

	}
	historyIfStat = newIfStat
	return
}
