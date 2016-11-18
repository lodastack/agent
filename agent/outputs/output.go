package outputs

import (
	"time"

	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
)

type OutputInf interface {
	// SetServers sets backend servers
	SetServers(servers []string)
	// Write takes in group of points to be written to the Output
	Write(chan Data)
	// Description returns a one-sentence description on the Output
	Description() string
}

type Creator func() OutputInf

var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}

type Data struct {
	Namespace string
	Points    *common.Points
}

var (
	Counter uint64
	queue   chan Data
)

func SendMetrics(ctype string, namespace string, metrics []*common.Metric) error {
	if len(metrics) == 0 {
		return nil
	}
	// filter topic
	namespace = "collect." + namespace

	data := &common.Points{Database: namespace, RetentionPolicy: "default", Precision: "s"}
	now := time.Now().Unix()
	hostname, err := common.Hostname()
	if err != nil {
		log.Errorf("get hostname failed: %s", err.Error())
		return err
	}
	for _, metric := range metrics {
		if metric.Tags == nil {
			metric.Tags = map[string]string{"host": hostname}
		} else {
			metric.Tags["host"] = hostname
		}
		if metric.Timestamp < 1e9 || metric.Timestamp > 1e10 {
			metric.Timestamp = now
		}
		log.Info("namespace:", namespace, " metric:", metric.String())
		p := &common.Point{metric.Name, metric.Timestamp, metric.Tags, map[string]interface{}{"value": metric.Value}}
		if ctype == common.TYPE_LOG {
			p.Fields["offset"] = metric.Offset
		}
		data.Points = append(data.Points, p)
	}

	queue <- Data{namespace, data}
	return nil
}

// Output runs collects data based on the given config.
type Output struct {
	Config *Config
}

// New returns an Output struct based off the given Config.
func New(config *Config) (*Output, error) {
	o := &Output{
		Config: config,
	}
	if o.Config.BufferSize <= 0 {
		o.Config.BufferSize = 1 << 16
	}
	return o, nil
}

func (o *Output) Start() {
	queue = make(chan Data, o.Config.BufferSize)
	creator, ok := Outputs[o.Config.Name]
	if !ok {
		panic("no output found: " + o.Config.Name)
	}
	output := creator()
	output.SetServers(o.Config.Servers)
	output.Write(queue)
}
