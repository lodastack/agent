package sysinfo

import (
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/outputs"

	"github.com/lodastack/log"
)

const updateIntrval = 10 * time.Second

func init() {
	go func() {
		for {
			if err := UpdateCpuStat(); err != nil {
				log.Errorf("update CPU status error: %s", err.Error())
			}
			if err := UpdateDiskStats(); err != nil {
				log.Errorf("update disks status error: %s", err.Error())
			}
			time.Sleep(updateIntrval)
		}
	}()
}

type Collector struct {
	Name  string
	Cycle int
	Fns   []func() []*common.Metric
}

func (self Collector) Run() {
	m := []*common.Metric{}
	for _, fn := range self.Fns {
		m = append(m, fn()...)
	}

	for _, ns := range common.GetNamespaces() {
		outputs.SendMetrics(self.Name, ns, m)
	}
}

func (self Collector) Description() string {
	return self.Name
}

func toMetric(name string, value interface{}, tags map[string]string) *common.Metric {
	ret := common.Metric{Name: name, Value: value, Tags: map[string]string{}}
	for k, v := range tags {
		ret.Tags[k] = v
	}
	return &ret
}
