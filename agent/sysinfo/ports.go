package sysinfo

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/outputs"

	"github.com/lodastack/log"
)

type PortCollector struct{}

func (self PortCollector) Run() {
	m := map[string][]*common.Metric{}
	reportPorts := common.ReportPorts()
	for _, p := range reportPorts {
		if isListening(p.Port, p.Timeout) {
			m[p.Namespace] = append(m[p.Namespace], toMetric(common.TYPE_PORT+"."+p.Name, 1, nil))
		} else {
			m[p.Namespace] = append(m[p.Namespace], toMetric(common.TYPE_PORT+"."+p.Name, 0, nil))
		}
	}
	for ns, ms := range m {
		outputs.SendMetrics(common.TYPE_PORT, ns, ms)
	}
}
func (self PortCollector) Description() string {
	return "PortCollector"
}

func isListening(port string, timeoutStr string) bool {
	var conn net.Conn
	var err error
	var timeout int
	if timeout, err = strconv.Atoi(timeoutStr); err != nil {
		log.Errorf("convert port timeout to int failed: %s", err)
		return false
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)
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
