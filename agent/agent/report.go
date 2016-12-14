package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/config"

	"github.com/lodastack/log"
	"github.com/lodastack/models"
)

func (a *Agent) Report() {
	for range time.NewTicker(time.Minute * 10).C {
		a.report()
	}
}

func (a *Agent) report() {
	hostname, err := common.Hostname()
	if err != nil {
		log.Error("get hostname failed: ", err)
		return
	}

	data := models.Report{
		UUID:        common.GetUUID(),
		IPList:      common.GetIpList(),
		Ns:          common.GetNamespaces(),
		Version:     config.Version,
		NewHostname: hostname,
		AgentType:   "loda-agent",
		Update:      false,
	}

	changed, OldHostname := common.HostnameChanged()
	if changed {
		data.OldHostname = OldHostname
		data.Update = true
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal failed: ", data)
	} else {
		url := fmt.Sprintf("http://%s/api/v1/agent/report", a.Config.RegistryAddr)
		resp, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Error("report agent info failed: ", err)
		} else {
			log.Info("report agent info successfully")
			if resp.StatusCode == http.StatusOK {
				file := filepath.Join(a.Config.PluginsDir, ".hostname")
				if err := ioutil.WriteFile(file, []byte(data.NewHostname), 0644); err != nil {
					log.Error("write hostname cache file failed: ", err)
				}
			}
			resp.Body.Close()
		}
	}
}
