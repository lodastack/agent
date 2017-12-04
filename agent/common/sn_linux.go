package common

import (
	"io/ioutil"
	"strings"

	"github.com/lodastack/log"
)

const snFile = "/sys/class/dmi/id/product_serial"

// SN return serial number of this machine
// no need to update serialNumber
func SN() string {
	if serialNumber != "" {
		return serialNumber
	}

	if !Exists(snFile) {
		return ""
	}
	read, err := ioutil.ReadFile(snFile)
	if err != nil {
		log.Error("Read file failed: ", err)
		return ""
	}
	sn := strings.Replace(string(read), " ", "", -1)
	serialNumber = sn
	return sn
}
