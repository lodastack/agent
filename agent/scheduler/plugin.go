package scheduler

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/log"
)

const (
	pluginDisableDir = "plugin.disabled"
)

func PluginStatus() map[string]bool {
	mutex.Lock()
	defer mutex.Unlock()
	res := map[string]bool{}
	for k := range pluginSchedulers {
		res[k] = !pluginDisabled[k]
	}
	return res
}

func DisablePlugin(ns, repo string) error {
	mutex.Lock()
	defer mutex.Unlock()
	dir := path.Join(common.Conf.PluginsDir, pluginDisableDir)
	if !common.Exists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Error("fail to mkdir, dir:", dir, " err:", err)
			return err
		}
	}
	p := ns + "|" + strings.Split(repo, "/")[1]
	_, err := os.Create(path.Join(dir, p))
	if err != nil {
		return err
	}
	pluginDisabled[p] = true
	if s, ok := pluginSchedulers[p]; ok {
		s.stop()
	} else {
		log.Info("plugin does not exist")
	}
	return nil
}

func EnablePlugin(ns, repo string) error {
	mutex.Lock()
	defer mutex.Unlock()
	return enablePlugin(ns + "|" + strings.Split(repo, "/")[1])
}

func enablePlugin(p string) error {
	err := os.Remove(path.Join(common.Conf.PluginsDir, pluginDisableDir, p))
	if err == os.ErrNotExist {
		err = nil
	}
	if err != nil {
		return err
	}
	delete(pluginDisabled, p)
	if s, ok := pluginSchedulers[p]; ok {
		go s.run()
	} else {
		return errors.New("no such plugin")
	}
	return nil
}

func loadDisabledPlugin() {
	dir := path.Join(common.Conf.PluginsDir, pluginDisableDir)
	if !common.Exists(dir) {
		return
	}
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error("failed to ReadDir of ", dir)
		return
	}
	for _, fi := range fis {
		pluginDisabled[fi.Name()] = true
	}
}

func autoEnablePlugin() {
	dir := path.Join(common.Conf.PluginsDir, pluginDisableDir)
	for {
		time.Sleep(time.Hour)
		mutex.Lock()
		if fis, err := ioutil.ReadDir(dir); err == nil {
			for _, fi := range fis {
				if time.Now().Sub(fi.ModTime()) > time.Hour*24 {
					enablePlugin(fi.Name())
					log.Info("auto enable plugin " + fi.Name())
				}
			}
		}
		mutex.Unlock()
	}
}
