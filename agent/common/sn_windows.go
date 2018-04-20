package common

import (
	"os/exec"
	"strings"

	"github.com/lodastack/log"
)

const snFile = "wmic bios get serialnumber"

// SN return serial number of this machine
// no need to update serialNumber
func SN() string {
	if serialNumber != "" {
		return serialNumber
	}

	parts := strings.Fields(snFile)
	head := parts[0]
	parts = parts[1:]

	cmd := exec.Command(head, parts...)
	out, err := cmd.Output()
	if err != nil {
		log.Errorf("get sn failed: %s", err.Error())
		return ""
	}
	list := strings.Split(string(out), "SerialNumber")
	serialNumber = strings.Replace(list[len(list)-1], " ", "", -1)
	return serialNumber
}
