package goplugin

import (
	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/outputs"

	"github.com/lodastack/log"
)

const (
	GOPLUGIN_RESULT = "goplugin.result"
)

type Collector struct {
	Namespace string                 `json:"Namespace"`
	Name      string                 `json:"name"`
	Interval  int                    `json:"interval"`
	Params    map[string]interface{} `json:"params"`
}

func (self Collector) Description() string {
	return self.Name
}

func (self Collector) Run() {
	f, ok := funcs[self.Name]
	if !ok {
		log.Error("there is no goplugin named ", self.Name)
		self.SubmitException()
		return
	}
	ms, err := f(self.Params)
	if err != nil {
		log.Error("failed to excute plugin ", self.Name, " err:", err)
		self.SubmitException()
		return
	}
	for _, m := range ms {
		m.Name = "PLUGIN." + self.Name + "." + m.Name
	}
	outputs.SendMetrics(common.TYPE_GOPLUGIN, self.Name, ms)
}

func (self Collector) SubmitException() {
	tags := map[string]string{"ns": self.Namespace, "name": self.Name}
	m := &common.Metric{Name: GOPLUGIN_RESULT, Tags: tags, Value: 1}
	outputs.SendMetrics(common.TYPE_GOPLUGIN, common.EXCEPTION_NS, []*common.Metric{m})
}
