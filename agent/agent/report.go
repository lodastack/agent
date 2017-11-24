package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/config"

	"github.com/lodastack/log"
	"github.com/lodastack/models"
)

// unit: minute
// agent report interval
const reportInterval = 10

func (a *Agent) Report() {
	a.report()
	// LB report API
	randNum := rand.Intn(reportInterval * 60 * 1000)
	time.Sleep(time.Duration(randNum) * time.Millisecond)

	for range time.NewTicker(time.Minute * reportInterval).C {
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
		SN:          common.SN(),
		NewIPList:   common.GetIpList(),
		Ns:          common.GetNamespaces(),
		Version:     config.Version,
		Commit:      config.Commit,
		Branch:      config.Branch,
		BuildTime:   config.BuildTime,
		GoVersion:   runtime.Version(),
		NewHostname: hostname,
		AgentType:   "loda-agent",
		Update:      false,
		UpdateTime:  time.Now(),
	}

	HostnameChanged, OldHostname := common.HostnameChanged()
	IPChanged, OldIP := common.IPChanged()
	if HostnameChanged || IPChanged {
		data.Update = true
	}
	data.OldHostname = OldHostname
	data.OldIPList = OldIP

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal failed: ", data)
	} else {
		url := fmt.Sprintf("%s/api/v1/agent/report", a.Config.RegistryAddr)
		resp, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Error("report agent info failed: ", err)
		} else {
			if resp.StatusCode == http.StatusOK {
				log.Info("report agent info successfully")
				file := filepath.Join(a.Config.PluginsDir, ".hostname")
				if err := ioutil.WriteFile(file, []byte(data.NewHostname), 0644); err != nil {
					log.Error("write hostname cache file failed: ", err)
				}
				file = filepath.Join(a.Config.PluginsDir, ".ip")
				if err := ioutil.WriteFile(file, []byte(strings.Join(data.NewIPList, ",")), 0644); err != nil {
					log.Error("write ip cache file failed: ", err)
				}
			} else {
				log.Errorf("report agent info failed: StatusCode %d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	}
}
