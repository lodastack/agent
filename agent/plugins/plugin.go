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
	Name        string   `json:"name"`
	Cycle       int      `json:"interval"`
	Namespace   string   `json:"namespace"`
	Repo        string   `json:"git"`
	ProjectName string   `json:projectname`
	Param       []string `json:"param"`
	Parameters  string   `json:"parameters"`
}

func (self Collector) Run() {
	err := self.Execute(self.Cycle*1000 - 500)
	if err != nil {
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

	bashFile := path.Join(dir, "plugin.sh")
	if !common.Exists(bashFile) {
		log.Error("failed to exec plugin ", self.Namespace+"|"+self.ProjectName, " plugin.sh doesn't exist")
		return errors.New("plugin.sh doesn't exist")
	}
	cmd := exec.Command("sh", append([]string{bashFile}, self.Param...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = dir
	err = cmd.Start()
	if err != nil {
		log.Error("fail to start plugin:", bashFile, " err:", err)
		return err
	}
	err, isTimeout := common.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Millisecond)
	if isTimeout {
		// has be killed
		if err == nil {
			log.Warning("timeout and kill process ", bashFile, " successfully")
		} else {
			log.Error("kill process ", bashFile, " occur error:", err)
		}
		return errors.New("plugin timeout")
	}

	if err != nil {
		log.Error("exec plugin ", bashFile, " failed. error:", err)
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
		log.Error("stdout of ", bashFile, " is blank")
		return nil
	}

	//For SecTeam
	namearr := strings.SplitN(self.ProjectName, "-", 2)
	if len(namearr) > 1 && namearr[0] == "SEC" {
		var metrics []*common.Metric
		err = json.Unmarshal(data, &metrics)
		if err != nil {
			log.Error("json.Unmarshal stdout of ", bashFile, " fail. error:", err, " stdout:", stdout.String())
			return err
		}

		for _, m := range metrics {
			m.Name = "SEC-PLUGIN." + self.Name + "." + m.Name
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
		log.Error("json.Unmarshal stdout of ", bashFile, " fail. error:", err, " stdout:", stdout.String())
		return err
	}

	for _, m := range metrics {
		m.Name = "PLUGIN." + self.Name + "." + m.Name
	}
	outputs.SendMetrics(common.TYPE_PLUGIN, self.Namespace, metrics)
	return nil
}
