package scheduler

import (
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/goplugin"
	"github.com/lodastack/agent/agent/loda"
	"github.com/lodastack/agent/agent/plugins"
	"github.com/lodastack/agent/agent/sysinfo"
	"github.com/lodastack/log"
)

// unit: second
// update NS interval
const updateInterval = 60

var (
	mutex sync.Mutex
)

var (
	pluginSchedulers map[string]*Scheduler
	pluginDisabled   map[string]bool
	sysSchedulers    map[string]*Scheduler

	GoPluginSchedulers map[string]*Scheduler //ns|name
)

func Start() {
	pluginSchedulers = make(map[string]*Scheduler)
	pluginDisabled = make(map[string]bool)
	sysSchedulers = make(map[string]*Scheduler)
	loadDisabledPlugin()

	GoPluginSchedulers = make(map[string]*Scheduler)

	go func() {
		Update()
		// LB update NS API
		randNum := rand.Intn(updateInterval * 1000)
		time.Sleep(time.Duration(randNum) * time.Millisecond)
		for {
			Update()
			time.Sleep(time.Second * updateInterval)
		}
	}()

	autoEnablePlugin()
}

func Update() {
	if namespaces, err := loda.Ns(); err == nil {
		common.SetNamespaces(namespaces)
		log.Info("get ns from loda: ", common.Namespaces)
	} else {
		log.Error("get ns from loda failed: ", err)
	}
	if newPorts, newProcs, newPlugins, newGoPlugins, intervals, err := loda.MonitorItems(); err == nil {
		updatePlugin(newPlugins)
		updateGoPlugin(newGoPlugins)
		updateSys(intervals)
		common.SetPorts(newPorts)
		common.SetProcs(newProcs)
	} else {
		log.Error("get monitor items from loda failed: ", err)
	}
}

func updatePlugin(collectors map[string]plugins.Collector) {
	mutex.Lock()
	defer mutex.Unlock()
	//plugin: ns|repo
	for plugin, scheduler := range pluginSchedulers {
		collector, ok := collectors[plugin]
		if !ok || !reflect.DeepEqual(collector, scheduler.collector) || collector.Cycle == 0 {
			scheduler.stop()
			delete(pluginSchedulers, plugin)
			log.Info("delete plugin:", plugin)
		}
	}
	for plugin, collector := range collectors {
		err := plugins.Update(collector.Namespace, common.GitPath(collector.Repo), false)
		if err != nil {
			log.Error("failed to update local plugin:", plugin, " from remote repo:", common.GitPath(collector.Repo), " err:", err)
			continue
		}
		if collector.Cycle == 0 {
			continue
		}
		if _, ok := pluginSchedulers[plugin]; ok {
			log.Info("plugin:", plugin, " already started")
			continue
		}
		s := NewScheduler(collector.Cycle, collector)
		pluginSchedulers[plugin] = s
		if !pluginDisabled[plugin] {
			go s.run()
		}
		log.Info("add plugin:", plugin)
	}
}

func updateGoPlugin(collectors map[string]goplugin.Collector) {
	mutex.Lock()
	defer mutex.Unlock()
	for name, s := range GoPluginSchedulers {
		newclt, ok := collectors[name]
		if !ok || !reflect.DeepEqual(newclt, s.collector) {
			log.Info("delete goplugin ", name)
			s.stop()
			delete(GoPluginSchedulers, name)
		}
	}
	for name, newclt := range collectors {
		if _, ok := GoPluginSchedulers[name]; ok {
			continue
		}
		s := NewScheduler(newclt.Interval, newclt)
		GoPluginSchedulers[name] = s
		go s.run()
		log.Info("add new goplugin:", name)
	}
}

func getFuncsByType(t string) (ret []func() []*common.Metric) {
	switch t {
	case common.TYPE_CPU:
		ret = append(ret, sysinfo.AgentMetrics, sysinfo.CpuMetrics, sysinfo.PsMetrics)
	case common.TYPE_DISK:
		ret = append(ret, sysinfo.IOStatsMetrics)
	case common.TYPE_MEM:
		ret = append(ret, sysinfo.MemMetrics)
	case common.TYPE_FS:
		ret = append(ret, sysinfo.FsKernelMetrics, sysinfo.FsRWMetrics,
			sysinfo.FsSpaceMetrics)
	case common.TYPE_KERNEL:
		ret = append(ret, sysinfo.WtmpMetrics, sysinfo.BtmpMetrics)
	case common.TYPE_TIME:
		ret = append(ret, sysinfo.TimeMetrics)
	case common.TYPE_DEV:
		ret = append(ret, sysinfo.PcapMetrics)
	case common.TYPE_NET:
		ret = append(ret, sysinfo.NetMetrics, sysinfo.SocketStatSummaryMetrics)
	case common.TYPE_COREDUMP:
		ret = append(ret, sysinfo.CoreDumpMetrics)
	}
	return ret
}

func updateSys(intervals map[string]int) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, t := range append(common.SYS_TYPES, common.TYPE_PORT, common.TYPE_PROC) {
		interval := intervals[t]
		s := sysSchedulers[t]
		if s == nil {
			if t == common.TYPE_PORT {
				s = NewScheduler(interval, sysinfo.PortCollector{})
			} else if t == common.TYPE_PROC {
				s = NewScheduler(interval, sysinfo.ProcCollector{})
			} else {
				s = NewScheduler(interval, sysinfo.Collector{t, interval, getFuncsByType(t)})
			}
			sysSchedulers[t] = s
			go s.run()
		} else if s.interval != interval {
			s.stop()
			s.setTicker(interval)
			go s.run()
		}
	}
}

func DeleteAll() {
	mutex.Lock()
	defer mutex.Unlock()
	// clean plugin
	for plugin, scheduler := range pluginSchedulers {
		scheduler.stop()
		delete(pluginSchedulers, plugin)
		log.Info("delete plugin:", plugin)
	}
	// clean sys
	for _, t := range append(common.SYS_TYPES, common.TYPE_PORT, common.TYPE_PROC) {
		if scheduler, ok := sysSchedulers[t]; ok {
			scheduler.stop()
			delete(sysSchedulers, t)
			log.Info("delete sys collect:", t)
		}
	}
}
