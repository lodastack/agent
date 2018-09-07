package plugins

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"os/user"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/outputs"

	"github.com/lodastack/log"
)

const (
	PLUGIN_RESULT = "plugin.result"
)

type Collector struct {
	Name        string `json:"name"`
	Cycle       int
	Namespace   string   `json:"namespace"`
	Repo        string   `json:"git"`
	ProjectName string   `json:projectname`
	Param       []string `json:"param"`
	Parameters  string   `json:"parameters"`
	StrCycle    string   `json:"interval"`
}

func (self Collector) Run() {
	err := self.Execute(self.Cycle*1000 - 500)
	if err != nil {
		log.Error("plugin execute failed: ", err)
		self.SubmitException()
	}
}

func (self Collector) SubmitException() {
	tags := map[string]string{"ns": self.Namespace, "name": self.Repo}
	if self.Cycle == 0 {
		tags["type"] = "call"
	} else {
		tags["type"] = "regular"
	}
	m := &common.Metric{Name: PLUGIN_RESULT, Tags: tags, Value: 1}
	outputs.SendMetrics(common.TYPE_SYS, common.EXCEPTION_NS, []*common.Metric{m})
}

func (self Collector) Description() string {
	return self.Namespace + "|" + self.ProjectName
}

func (self Collector) Execute(timeout int) error {
	dir := path.Join(common.Conf.PluginsDir, self.Namespace, self.ProjectName)
	execUser, err := user.Lookup(common.Conf.PluginsUser)
	if err != nil {
		log.Error("can not su to plugin user: ", err)
		return err
	}
	uid, err := strconv.Atoi(execUser.Uid)
	if err != nil {
		log.Error("failed to get uid of plugin user: ", err)
		return err
	}
	gid, err := strconv.Atoi(execUser.Gid)
	if err != nil {
		log.Error("failed to get gid of plugin user: ", err)
		return err
	}

	pluginFile := path.Join(dir, "plugin")
	if !common.Exists(pluginFile) {
		log.Error("failed to exec plugin ", self.Namespace+"|"+self.ProjectName, " plugin doesn't exist")
		return errors.New("plugin doesn't exist")
	}
	cmd := exec.Command(pluginFile, self.Param...)
	//Create session
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = dir
	err = cmd.Start()
	if err != nil {
		log.Error("fail to start plugin:", pluginFile, " err:", err)
		return err
	}
	err, isTimeout := common.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Millisecond)
	if isTimeout {
		// has be killed
		if err == nil {
			log.Warning("timeout and kill process ", pluginFile, " successfully")
		} else {
			log.Error("kill process ", pluginFile, " occur error:", err)
		}
		return errors.New("plugin timeout")
	}

	if err != nil {
		log.Error("exec plugin ", pluginFile, " failed. error:", err)
		log.Debug("stdout: ", string(stdout.Bytes()))
		log.Debug("stderr: ", string(stderr.Bytes()))
		return err
	}
	if self.Cycle == 0 {
		return nil
	}

	// exec successfully
	data := stdout.Bytes()
	if len(data) == 0 {
		log.Error("stdout of ", pluginFile, " is blank")
		return nil
	}

	//For SecTeam
	namearr := strings.SplitN(self.ProjectName, "-", 2)
	if len(namearr) > 1 && namearr[0] == "SEC" {
		var metrics []*common.Metric
		err = json.Unmarshal(data, &metrics)
		if err != nil {
			log.Error("json.Unmarshal stdout of ", pluginFile, " fail. error:", err, " stdout:", stdout.String())
			return err
		}

		for _, m := range metrics {
			m.Name = "SEC." + self.Name + "." + m.Name
		}
		nsarr := strings.SplitN(self.Namespace, ".", 2)
		if len(nsarr) > 1 {
			outputs.SendMetrics(common.TYPE_PLUGIN, "sec."+nsarr[1], metrics)
		} else {
			log.Warning("SEC plugin repo namespace error")
		}
		return nil
	}

	var metrics []*common.Metric
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		log.Error("json.Unmarshal stdout of ", pluginFile, " fail. error:", err, " stdout:", stdout.String())
		return err
	}

	for index, m := range metrics {
		var customNS string
		if ns, ok := m.Tags["namespace"]; ok {
			ns = strings.Replace(ns, "-", ".", -1)
			if strings.HasSuffix(ns, ".loda") {
				customNS = ns
			}
		}
		if err := nameCheck(m.Name); err != nil {
			metrics = append(metrics[:index], metrics[index+1:]...)
			log.Errorf("metrics check failed: %s", err)
			continue
		}
		m.Name = self.Name + "." + m.Name
		if customNS != "" {
			outputs.SendMetrics(common.TYPE_PLUGIN, customNS, []*common.Metric{m})
		}
	}
	outputs.SendMetrics(common.TYPE_PLUGIN, self.Namespace, metrics)
	return nil
}

func nameCheck(name string) error {
	for _, nameLetter := range name {
		if nameLetter == '-' || nameLetter == '_' || nameLetter == '.' || (nameLetter >= 'a' && nameLetter <= 'z') || (nameLetter >= 'A' && nameLetter <= 'Z') || (nameLetter >= '0' && nameLetter <= '9') {
			continue
		}
		return errors.New("invalid metric name, just allow 0-9 a-z A-Z - _ .")
	}
	return nil
}
