package agent

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
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

	if !common.Exists(a.Config.PluginsDir) {
		if err := os.MkdirAll(a.Config.PluginsDir, 0755); err != nil {
			log.Error("create hostname cache dir failed: ", err)
			return
		}
	}
	file := filepath.Join(a.Config.PluginsDir, ".hostname")
	//read saved content
	read, err := ioutil.ReadFile(file)
	if os.IsNotExist(err) {
		if err := ioutil.WriteFile(file, []byte(data.NewHostname), 0644); err != nil {
			log.Error("write hostname cache file failed: ", err)
		}
	}
	if err != nil {
		log.Error("Read hostname cache file failed: ", err)
	}

	if err == nil {
		if string(read) != data.NewHostname {
			log.Infof("Hostname chaged: %s -> %s", string(read), data.NewHostname)
			data.OldHostname = string(read)
			data.Update = true
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal failed: ", data)
	} else {
		resp, err := http.Post(a.Config.ReportAddr, "application/json;charset=utf-8", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Error("report agent info failed: ", err)
		} else {
			log.Info("report agent info successfully")
			if resp.StatusCode == http.StatusOK {
				if err := ioutil.WriteFile(file, []byte(data.NewHostname), 0644); err != nil {
					log.Error("write hostname cache file failed: ", err)
				}
			}
			resp.Body.Close()
		}
	}
}
