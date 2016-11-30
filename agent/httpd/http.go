package httpd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/outputs"
	"github.com/lodastack/agent/agent/plugins"
	"github.com/lodastack/agent/agent/scheduler"
	"github.com/lodastack/agent/config"

	"github.com/lodastack/log"
)

var (
	runnningPlugins = make(map[string]bool)
	mutex           sync.Mutex
)

func PluginListHandler(w http.ResponseWriter, req *http.Request) {
	p := scheduler.PluginStatus()
	res := ""
	for k, v := range p {
		res += fmt.Sprintf("%s enabled:%v\n", k, v)
	}
	io.WriteString(w, res)
}

func PluginDisableHandler(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	namespace := q.Get("ns")
	repo := q.Get("repo")
	repo, err := checkRepo(namespace, repo)
	if err != nil {
		io.WriteString(w, err.Error()+"\n")
		return
	}
	err = scheduler.DisablePlugin(namespace, repo)
	if err != nil {
		io.WriteString(w, "failed to disable plugin: "+err.Error()+"\n")
	} else {
		io.WriteString(w, "plugin disabled successfully\n")
	}
}

func PluginEnableHandler(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	namespace := q.Get("ns")
	repo := q.Get("repo")
	repo, err := checkRepo(namespace, repo)
	if err != nil {
		io.WriteString(w, err.Error()+"\n")
		return
	}
	err = scheduler.EnablePlugin(namespace, repo)
	if err != nil {
		io.WriteString(w, "failed to Enable plugin: "+err.Error()+"\n")
	} else {
		io.WriteString(w, "plugin enabled successfully\n")
	}
}

func PluginUpdateHandler(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	namespace := q.Get("ns")
	repo := q.Get("repo")
	repo, err := checkRepo(namespace, repo)
	if err != nil {
		io.WriteString(w, err.Error()+"\n")
		return
	}
	err = plugins.Update(namespace, common.GitPath(repo), true)
	if err != nil {
		io.WriteString(w, "update failed: "+err.Error()+"\n")
	} else {
		io.WriteString(w, "update successfully\n")
	}
}

func PluginRunHandler(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	namespace := q.Get("ns")
	repo := q.Get("repo")
	timeouts := q.Get("timeout")
	parameters := q.Get("param")
	param := strings.Split(parameters, " ")
	log.Info("recieve a plugin/run request: ", namespace, " ", repo, " ", param)
	repo, err := checkRepo(namespace, repo)
	if err != nil {
		io.WriteString(w, err.Error()+"\n")
		return
	}
	var timeout int
	if timeouts != "" {
		t, err := strconv.ParseUint(timeouts, 10, 64)
		if err != nil {
			io.WriteString(w, "invalid timeout\n")
			return
		}
		timeout = int(t)
	} else {
		timeout = 10
	}
	//check this plugin exist
	err = plugins.Update(namespace, common.GitPath(repo), false)
	if err != nil {
		io.WriteString(w, "update failed: "+err.Error()+"\n")
		return
	}
	mutex.Lock()
	plugin := strings.Split(repo, "/")[1]
	s := namespace + "|" + plugin
	if runnningPlugins[s] {
		io.WriteString(w, "plugin is running")
		mutex.Unlock()
		return
	} else {
		runnningPlugins[s] = true
		mutex.Unlock()
	}

	c := plugins.Collector{"", 0, namespace, repo, plugin, param, parameters, "0"}
	err = c.Execute(timeout * 1000)
	if err != nil {
		c.SubmitException()
		io.WriteString(w, "run plugin "+namespace+" "+plugin+" accurs an error:"+err.Error()+"\n")
	} else {
		io.WriteString(w, "run plugin successfully")
	}
	delete(runnningPlugins, s)
}

func PostDataHandler(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	namespace := q.Get("ns")
	if namespace == "" {
		io.WriteString(w, "invalid ns\n")
		return
	}
	decoder := json.NewDecoder(req.Body)
	var m []*common.Metric
	err := decoder.Decode(&m)
	if err != nil {
		io.WriteString(w, "failed to decode json:"+err.Error()+"\n")
		return
	}
	outputs.SendMetrics(common.TYPE_SYS, namespace, m)
	io.WriteString(w, "send data to nsq\n")
}

func UpdateHandlder(w http.ResponseWriter, req *http.Request) {
	scheduler.Update()
	io.WriteString(w, "collect items updated\n")
}

func GetNsHandler(w http.ResponseWriter, req *http.Request) {
	ns := common.GetNamespaces()
	rep, err := json.Marshal(ns)
	if err != nil {
		io.WriteString(w, "[]")
	}
	io.WriteString(w, string(rep))
}

func GetStatusHandler(w http.ResponseWriter, req *http.Request) {
	status := make(map[string]string)
	status["version"] = config.Version
	rep, err := json.Marshal(status)
	if err != nil {
		io.WriteString(w, "{}")
	}
	io.WriteString(w, string(rep))
}

func LogOffsetHandler(w http.ResponseWriter, req *http.Request) {
	var result common.Result
	q := req.URL.Query()
	fpath := q.Get("fpath")
	offset, err := strconv.ParseInt(q.Get("offset"), 10, 64)
	if err != nil {
		result.StatusCode = 500
		result.Msg = "parse param:offset error"
		io.WriteString(w, result.String())
		return
	}
	lineNum, err := strconv.ParseInt(q.Get("num"), 10, 64)
	if err != nil {
		result.StatusCode = 500
		result.Msg = "parse param:num error"
		io.WriteString(w, result.String())
		return
	}
	lines, err := common.ReadLinesFromOffset(fpath, offset, lineNum)
	if err != nil {
		result.StatusCode = 500
		result.Msg = fmt.Sprintf("read lines from offset failed. err:%s", err.Error())
		result.Data = lines
		io.WriteString(w, result.String())
		return
	}
	result.StatusCode = 200
	result.Msg = fmt.Sprintf("read %d lines of log start from offset:%d successfully", lineNum, offset)
	result.Data = lines
	io.WriteString(w, result.String())
}

func nameFromGit(repo string) (string, error) {
	if strings.Count(repo, ":") != 1 || !strings.HasSuffix(repo, ".git") {
		return "", errors.New("invalid git repo path")
	}
	repo = strings.Split(repo, ":")[1]
	repo = repo[:len(repo)-4]
	if strings.Count(repo, "/") != 1 {
		return "", errors.New("invalid git path")
	}
	return repo, nil
}

func checkPlugin(ns, repo string) bool {
	return common.GetPluginInfo()[ns+"|"+repo]
}

func checkRepo(namespace, repo string) (string, error) {
	if namespace == "" || repo == "" {
		return "", errors.New("invalid ns or repo")
	}
	repo, err := nameFromGit(repo)
	if err != nil {
		return "", errors.New("invalid repo")
	}
	if !checkPlugin(namespace, repo) {
		return "", errors.New("plugin doesn't exist, add to odin first please")
	}
	return repo, nil
}

func Start(listen string) {
	http.HandleFunc("/plugins/list", PluginListHandler)
	http.HandleFunc("/plugins/update", PluginUpdateHandler)
	http.HandleFunc("/plugins/run", PluginRunHandler)
	http.HandleFunc("/plugins/disable", PluginDisableHandler)
	http.HandleFunc("/plugins/enable", PluginEnableHandler)
	http.HandleFunc("/post", PostDataHandler)
	http.HandleFunc("/update", UpdateHandlder)
	http.HandleFunc("/me/ns", GetNsHandler)
	http.HandleFunc("/me/status", GetStatusHandler)
	//http.HandleFunc("/log/offset", LogOffsetHandler)
	fmt.Println("starting collect module http listener... on ", common.Conf.Listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}
