package sysinfo

import (
	"errors"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

func FsKernelMetrics() (L []*common.Metric) {
	maxFiles, err := nux.KernelMaxFiles()
	if err != nil {
		log.Error("failed collect kernel metrics:", err)
		return
	}

	L = append(L, toMetric("fs.files.max", maxFiles, nil))

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Error("failed to call KernelAllocateFiles:", err)
		return
	}

	v := math.Ceil(float64(allocateFiles) * 100 / float64(maxFiles))
	L = append(L, toMetric("fs.files.allocated", allocateFiles, nil))
	L = append(L, toMetric("fs.files.used.percent", v, nil))
	L = append(L, toMetric("fs.files.left", maxFiles-allocateFiles, nil))
	return
}

func FsRWMetrics() (L []*common.Metric) {
	mountPoints, err := nux.ListMountPoint()

	if err != nil {
		log.Error("failed to call ListMountPoint:", err)
		return
	}

	for idx := range mountPoints {
		var du *nux.DeviceUsage
		var res int
		du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
		if err != nil {
			log.Error("failed to call BuildDeviceUsage:", err)
			continue
		}
		file := filepath.Join(du.FsFile, ".loda-fs-detect")
		now := time.Now().Format("2006-01-02 15:04:05")
		content := "FS-RW" + now
		err = CheckFS(file, content)
		if err != nil {
			res = 0
		} else {
			res = 1
		}
		tags := map[string]string{"mount": du.FsFile}
		L = append(L, toMetric("fs.disk.rw", res, tags))
	}

	return
}

func CheckFS(file string, content string) error {
	//write test
	fd, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	defer fd.Close()
	if err != nil {
		log.Error("Open file failed: ", err)
		return err
	}
	buf := []byte(content)
	count, err := fd.Write(buf)
	if err != nil || count != len(buf) {
		log.Error("Write file failed: ", err)
		return err
	}
	//read test
	read, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error("Read file failed: ", err)
		return err
	}
	if string(read) != content {
		log.Error("Read content failed: ", string(read))
		return errors.New("Read content failed")
	}
	//clean the file
	err = os.Remove(file)
	if err != nil {
		log.Error("Remove file filed: ", err)
		return err
	}
	return nil
}
