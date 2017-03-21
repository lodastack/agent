package loda

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/goplugin"
	"github.com/lodastack/agent/agent/plugins"
	"github.com/lodastack/log"
)

var Zerotimes = 0

func Get(url string) (b []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Get url failed: %s  not found", url)
		return
	}
	b, err = ioutil.ReadAll(resp.Body)
	return
}

func Post(url string, data []byte) ([]byte, error) {
	body := bytes.NewBuffer([]byte(data))
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("Post url failed: %s code: %d", url, res.StatusCode)
		return nil, err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Ns() ([]string, error) {
	var res []string
	host, err := common.Hostname()
	if err != nil {
		return res, err
	}
	//check hostname chaged
	changed, _ := common.HostnameChanged()
	if changed {
		return res, fmt.Errorf("hostname changed, skip fetch ns")
	}

	url := fmt.Sprintf("%s/api/v1/agent/ns", common.Conf.RegistryAddr)
	data := make(map[string]string)
	data["hostname"] = host
	data["ip"] = strings.Join(common.GetIpList(), ",")
	byteData, err := json.Marshal(data)
	b, err := Post(url, byteData)
	if err != nil {
		return res, err
	}

	type ResponseNS struct {
		Code int               `json:"httpstatus"`
		Data map[string]string `json:"data"`
	}
	var response ResponseNS
	err = json.Unmarshal(b, &response)
	if err != nil {
		log.Warning("json.Marshal Ns failed: ", err)
		return res, err
	}

	resp := response.Data

	var ids []string
	for ns, id := range resp {
		res = append(res, ns)
		ids = append(ids, id)
	}
	common.SetUUID(ids)
	return res, nil
}

func pullResources(ns string) (res []map[string]string, err error) {
	url := fmt.Sprintf("%s/api/v1/agent/resource?ns=%s&type=collect", common.Conf.RegistryAddr, ns)
	b, err := Get(url)
	if err != nil {
		return
	}

	type ResponseRes struct {
		Code int                 `json:"httpstatus"`
		Data []map[string]string `json:"data"`
	}
	var response ResponseRes
	err = json.Unmarshal(b, &response)
	if err != nil {
		return
	}
	res = response.Data

	if len(res) == 0 {
		err = fmt.Errorf("no items under this namespace")
		Zerotimes++
	} else {
		Zerotimes = 0
	}
	return
}

func MonitorItems() (ports []common.PortMonitor,
	procs []common.ProcMonitor,
	pluginCollectors map[string]plugins.Collector,
	gopluginCollectors map[string]goplugin.Collector,
	intervals map[string]int, err error) {
	nss := common.GetNamespaces()
	pluginCollectors = make(map[string]plugins.Collector)
	pluginInfo := make(map[string]bool)
	gopluginCollectors = make(map[string]goplugin.Collector)
	intervals = make(map[string]int)

	for _, ns := range nss {
		err = nil
		var items []map[string]string
		items, err = pullResources(ns)
		if err != nil {
			log.Error("failed to get resources from registry, ns: ", ns, " err: ", err)
			continue
		}

		for _, item := range items {
			monitorType, ok := item["measurement_type"]
			if !ok {
				log.Warning("measurement_type is not exist: ", item["measurement_type"])
				continue
			}
			b, err := json.Marshal(item)
			if err != nil {
				log.Warning("json.Marshal item failed: ", err)
				continue
			}
			switch monitorType {
			case common.TYPE_PORT:
				if port, err := parsePort(b); err == nil {
					port.Namespace = ns
					ports = append(ports, port)
				} else {
					log.Warning("parsePort failed: ", err)
				}
			case common.TYPE_PROC:
				if proc, err := parseProc(b); err == nil {
					proc.Namespace = ns
					procs = append(procs, proc)
				} else {
					log.Warning("parseProc failed: ", err)
				}
			case common.TYPE_PLUGIN:
				if col, err := parsePlugin(b); err == nil {
					col.Namespace = ns
					pluginCollectors[ns+"|"+col.ProjectName] = col
					pluginInfo[ns+"|"+col.Repo] = true
				} else {
					log.Warning("get plugin collection failed: ", err)
				}
			case common.TYPE_GOPLUGIN:
				if col, err := parseGoPlugin(b); err == nil {
					col.Namespace = ns
					gopluginCollectors[ns+"|"+col.Name] = col
				} else {
					log.Warning("parse goplugin collection failed: ", err)
				}
			default:
				if interval, ok := item["interval"]; ok {
					var intervalInt int
					if intervalInt, err = strconv.Atoi(interval); err == nil {
						intervals[monitorType] = intervalInt
					} else {
						log.Warning("convert interval to int failed: ", err)
					}
				}
			}
		}
		//getAlarmPlugin(ns, pluginCollectors, pluginInfo)
	}

	common.SetPluginInfo(pluginInfo)
	for _, t := range append(common.SYS_TYPES, common.TYPE_PORT, common.TYPE_PROC) {
		if intervals[t] == 0 {
			intervals[t] = common.DEFAULT_INTERVAL[t]
		}
	}
	return
}

// func getAlarmPlugin(ns string, pluginCollectors map[string]plugins.Collector, pluginInfo map[string]bool) {
// 	url := fmt.Sprintf("http://%s/api/v1/resource?ns=%s&resource=alarm", common.Conf.RegistryAddr, ns)
// 	b, err := Get(url)
// 	if err != nil {
// 		return
// 	}
// 	var response models.Response
// 	err = json.Unmarshal(b, &response)
// 	if err != nil {
// 		log.Error("Unmarshal from alarm failed: ", err)
// 		return
// 	}

// 	alarms, ok := response.Data.([]map[string]interface{})
// 	if !ok {
// 		err = fmt.Errorf("response data is not a map slice type")
// 		return
// 	}

// 	for _, alarm := range alarms {
// 		ac, ok := alarm["actions"]
// 		if !ok {
// 			continue
// 		}
// 		actions, ok := ac.([]interface{})
// 		if !ok {
// 			continue
// 		}
// 		for _, a := range actions {
// 			action, ok := a.(map[string]interface{})
// 			if !ok {
// 				continue
// 			}
// 			t, ok := action["type"]
// 			if !ok {
// 				continue
// 			}
// 			if t1, ok := t.(string); ok && t1 == "AGENT" {
// 				b, err := json.Marshal(action)
// 				if err != nil {
// 					log.Error("json.Marshal action failed: ", err)
// 					continue
// 				}
// 				var col plugins.Collector
// 				err = json.Unmarshal(b, &col)
// 				if err != nil {
// 					log.Error("json.Unmarshal to plugins.Collector from alarm failed: ", err)
// 					continue
// 				}
// 				col, err = formatPlugin(col)
// 				if err != nil {
// 					continue
// 				}
// 				col.Cycle = 0
// 				pluginInfo[ns+"|"+col.Repo] = true
// 				//col.Namespace = ns
// 				//pluginCollectors[ns+"|"+col.ProjectName] = col
// 			}
// 		}
// 	}
// }

func parsePlugin(b []byte) (res plugins.Collector, err error) {
	if err = json.Unmarshal(b, &res); err != nil {
		return
	}
	var cycle int
	if cycle, err = strconv.Atoi(res.StrCycle); err != nil {
		return
	}
	res.Cycle = cycle
	res, err = formatPlugin(res)
	return
}

func formatPlugin(p plugins.Collector) (plugins.Collector, error) {
	if strings.Contains(p.Repo, ":") {
		s := strings.Split(p.Repo, ":")[1]
		p.Repo = s[:len(s)-4]
	}
	if strings.Count(p.Repo, "/") == 1 {
		p.ProjectName = strings.Split(p.Repo, "/")[1]
	}
	if p.Parameters != "" {
		for _, s := range strings.Split(p.Parameters, " ") {
			if s != "" {
				if strings.ContainsAny(s, ";|&<>`") {
					return p, errors.New("Invalid parameter")
				}
				p.Param = append(p.Param, s)
			}
		}
	}
	return p, nil
}

func parseProc(b []byte) (res common.ProcMonitor, err error) {
	err = json.Unmarshal(b, &res)
	return
}

func parsePort(b []byte) (res common.PortMonitor, err error) {
	err = json.Unmarshal(b, &res)
	return
}

func parseGoPlugin(b []byte) (res goplugin.Collector, err error) {
	err = json.Unmarshal(b, &res)
	return
}
