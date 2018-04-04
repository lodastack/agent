package common

import (
	"net"
	"strings"
)

// IP gets all IP from given if prefix
func IP() (ips []string, err error) {
	ips = make([]string, 0)

	ifaces, e := net.Interfaces()
	if e != nil {
		return ips, e
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		// ignore docker and warden bridge
		if !HasInterfacePrefix(iface.Name) {
			continue
		}

		addrs, e := iface.Addrs()
		if e != nil {
			return ips, e
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// IP filter
			// 224.0.0
			// 169.254.0.0/16
			// ff02::/16
			// ffx2::/16
			if ip == nil || ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
				continue
			}

			ipStr := ip.String()
			// append all IP
			ips = append(ips, ipStr)
		}
	}

	return ips, nil
}

func HasInterfacePrefix(ifacename string) bool {
	for _, prefix := range Conf.IfacePrefix {
		if strings.HasPrefix(ifacename, prefix) {
			return true
		}
	}
	return false
}
