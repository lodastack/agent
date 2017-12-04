package common

import (
	"fmt"
	"os/exec"
	"strings"
)

// SN return serial number of this machine
// no need to update serialNumber
func SN() string {
	if serialNumber != "" {
		return serialNumber
	}
	var sn string
	out, err := exec.Command("/usr/sbin/ioreg", "-l").Output()
	if err != nil {
		return ""
	}

	for _, l := range strings.Split(string(out), "\n") {
		if strings.Contains(l, "IOPlatformSerialNumber") {
			s := strings.Split(l, " ")
			sn = fmt.Sprintf("%s", s[len(s)-1])
			sn = strings.Replace(sn, "\"", "", -1)
		}
	}

	serialNumber = sn
	return sn
}
