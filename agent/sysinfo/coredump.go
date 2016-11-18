package sysinfo

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"github.com/lodastack/agent/agent/common"
	//"github.com/lodastack/agent/agent/outputs"
)

const (
	COREDUMP_DIR      = "/home/coresave"
	PATTERN           = "^core.(?P<service>[a-zA-Z0-9_-]+).(?P<pid>[0-9]+).(?P<timestamp>[0-9]+)$"
	COREDUMP_INTERVAL = 60
)

func CoreDumpMetrics() (L []*common.Metric) {
	name := "app.service.coredump"
	fis, err := ioutil.ReadDir(COREDUMP_DIR)
	if err != nil {
		//m := toMetric(name, 1, map[string]string{"service": "coredump-dir-not-found"})
		//outputs.SendMetrics(common.TYPE_SYS, common.EXCEPTION_NS, []*common.Metric{m})
	} else {
		ts := time.Now().Unix()
		re := regexp.MustCompile(PATTERN)
		for _, fi := range fis {
			values := re.FindStringSubmatch(fi.Name())
			if len(values) == 4 {
				service := values[1]
				timestamp, _ := strconv.ParseInt(values[3], 10, 64)
				if timestamp >= ts-COREDUMP_INTERVAL {
					L = append(L, toMetric(name, 1, map[string]string{"service": service}))
				}
			}
		}
	}
	return
}
