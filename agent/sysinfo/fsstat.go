package sysinfo

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

func FsSpaceMetrics() (L []*common.Metric) {
	mountPoints, err := nux.ListMountPoint()

	if err != nil {
		log.Error("failed to call ListMountPoint:", err)
		return
	}

	for idx := range mountPoints {
		var du *nux.DeviceUsage
		du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
		if err != nil {
			log.Error("failed to call BuildDeviceUsage:", err)
			continue
		}

		tags := map[string]string{"mount": du.FsFile}
		L = append(L, toMetric("fs.inodes.used.percent", du.InodesUsedPercent, tags))
		L = append(L, toMetric("fs.space.used.percent", du.BlocksUsedPercent, tags))
	}

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
		err = CheckFS(file, content, mountPoints[idx][3])
		if err != nil {
			res = 0
		} else {
			res = 1
		}
		tags := map[string]string{"mount": du.FsFile}
		L = append(L, toMetric("fs.files.rw", res, tags))
	}

	return
}

func CheckFS(file string, content string, t string) error {
	//  var t from /proc/mounts
	//  We can not check read only file system, we can not write a file
	if t == "rw" {
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
	}
	return nil
}
