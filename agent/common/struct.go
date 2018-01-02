package common

import (
	"fmt"
)

type Collector interface {
	Run() // collect metrics, call SendToNsq
	Description() string
}

type Metric struct {
	Name      string            `json:"name"`
	Timestamp int64             `json:"timestamp"`
	Tags      map[string]string `json:"tags"`
	Value     interface{}       `json:"value"`
	Offset    int64             `json:"offset,omitempty"`
}

// String series metric
func (m *Metric) String() string {
	//offset might by empty(0)
	return fmt.Sprintf("<%s %d %s %v %d>", m.Name, m.Timestamp, m.Tags, m.Value, m.Offset)
}

// key returns metric key
func (m *Metric) Key() string {
	return fmt.Sprintf("<%s%d%s>", m.Name, m.Timestamp, m.Tags)
}

type Point struct {
	Measurement string                 `json:"measurement"`
	Timestamp   int64                  `json:"timestamp"`
	Tags        map[string]string      `json:"tags"`
	Fields      map[string]interface{} `json:"fields"`
}

type Points struct {
	Database        string   `json:"database"`
	RetentionPolicy string   `json:"retentionPolicy"`
	Points          []*Point `json:"points"`
	Precision       string   `json:"precision"`
}

type PortMonitor struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Port      string `json:"port"`
	Timeout   string `json:"connect_timeout"`
}

type ProcMonitor struct {
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	BinaryPath string `json:"bin_path"`
}
