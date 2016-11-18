package sysinfo

import (
	"fmt"
	"net"
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/outputs"
)

type PortCollector struct{}

func (self PortCollector) Run() {
	m := map[string][]*common.Metric{}
	reportPorts := common.ReportPorts()
	for _, p := range reportPorts {
		if isListening(p.Port, p.Timeout) {
			m[p.Namespace] = append(m[p.Namespace], toMetric(p.Name, 1, nil))
		} else {
			m[p.Namespace] = append(m[p.Namespace], toMetric(p.Name, 0, nil))
		}
	}
	for ns, ms := range m {
		outputs.SendMetrics(common.TYPE_PORT, ns, ms)
	}
}
func (self PortCollector) Description() string {
	return "PortCollector"
}

func isListening(port int, timeout int) bool {
	var conn net.Conn
	var err error
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	if timeout <= 0 {
		// default timeout 3 second
		timeout = 3
	}
	conn, err = net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
