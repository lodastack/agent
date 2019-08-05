package common

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/lodastack/log"
)

func Hostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return normalizedHostname(hostname), nil
}

func normalizedHostname(hostname string) string {
	if strings.HasSuffix(hostname, HOST_SUFFIX) {
		return strings.TrimSuffix(hostname, HOST_SUFFIX)
	}
	return hostname
}

func GetIpList() []string {
	ips, err := IP()
	if err != nil {
		res := []string{}
		return res
	}
	return ips
}

func CmdRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):

		go func() {
			<-done // allow goroutine to exit
		}()

		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			err = syscall.Kill(-pgid, 15)
		} else {
			err = cmd.Process.Kill()
		}
		return err, true
	case err = <-done:
		return err, false
	}
}

func StrTagsToMap(tags string) map[string]string {
	res := map[string]string{}
	for _, tag := range strings.Split(tags, ",") {
		if len(tag) > 0 {
			pairs := strings.Split(tag, "=")
			if len(pairs) == 2 {
				res[pairs[0]] = pairs[1]
			}
		}
	}
	return res
}

func GitPath(repo string) string {
	return fmt.Sprintf(Conf.Git, repo)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func SetPrecision(from float64, precision int) float64 {
	base := math.Pow10(precision)
	return float64(int64(from*base)) / base
}

func ReadLinesFromOffset(fpath string, offset int64, lineNum int64) (lines []string, err error) {
	f, err := os.Open(fpath)
	defer f.Close()
	if err != nil {
		return
	}
	_, err = f.Seek(offset, 0)
	if err != nil {
		return
	}

	r := bufio.NewReader(f)
	var next int64
	for next = 0; next < lineNum; next++ {
		line, err := r.ReadString('\n')
		if err == nil {
			lines = append(lines, line)
		} else {
			break
		}
	}
	return
}

func HostnameChanged() (bool, string) {
	h, err := Hostname()
	if err != nil || h == "" {
		return false, ""
	}

	if !Exists(Conf.PluginsDir) {
		if err := os.MkdirAll(Conf.PluginsDir, 0755); err != nil {
			log.Error("create hostname cache dir failed: ", err)
			return false, h
		}
	}
	file := filepath.Join(Conf.PluginsDir, ".hostname")
	//read saved content
	read, err := ioutil.ReadFile(file)
	if os.IsNotExist(err) {
		if err := ioutil.WriteFile(file, []byte(h), 0644); err != nil {
			log.Error("write hostname cache file failed: ", err)
			return false, h
		}
	}
	if err != nil {
		log.Error("Read hostname cache file failed: ", err)
		return false, h
	}

	if strings.TrimSpace(string(read)) == "" {
		ioutil.WriteFile(file, []byte(h), 0644)
		return false, h
	}

	if string(read) != h {
		log.Infof("Hostname chaged: %s -> %s", string(read), h)
		return true, string(read)
	}
	return false, h
}

func IPChanged() (bool, []string) {
	var oldIP []string
	newIP := GetIpList()
	newIPContent := strings.Join(newIP, ",")

	if !Exists(Conf.PluginsDir) {
		if err := os.MkdirAll(Conf.PluginsDir, 0755); err != nil {
			log.Error("create IP cache dir failed: ", err)
			return false, oldIP
		}
	}
	file := filepath.Join(Conf.PluginsDir, ".ip")
	//read saved content
	read, err := ioutil.ReadFile(file)
	if os.IsNotExist(err) {
		if err := ioutil.WriteFile(file, []byte(newIPContent), 0644); err != nil {
			log.Error("write IP cache file failed: ", err)
			return false, oldIP
		}
	}
	if err != nil {
		log.Error("Read IP cache file failed: ", err)
		return false, oldIP
	}
	oldIP = strings.Split(string(read), ",")
	if string(read) != newIPContent {
		log.Infof("IP chaged: %s -> %s", string(read), newIPContent)
		return true, oldIP
	}
	return false, oldIP
}
