package sysinfo

import (
	"net"
	"sync"
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// PcapMetrics dep pcap lib
func PcapMetrics() (L []*common.Metric) {
	var mu sync.Mutex
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Error("collect net interfaces error:", err)
		return
	}
	var wg sync.WaitGroup
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		// ignore docker and warden bridge
		if !common.HasInterfacePrefix(iface.Name) {
			continue
		}

		// Get IP address
		var IPAddress string
		addrs, e := iface.Addrs()
		if e != nil {
			continue
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
			if ip == nil || ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			IPAddress = ip.String()
			break
		}

		if IPAddress == "" {
			continue
		}

		wg.Add(1)
		go func(ifname string, ifip string) {
			defer wg.Add(-1)
			// Open device
			handle, err := pcap.OpenLive(ifname, snapshotLen, promiscuous, pcaptimeout)
			if err != nil {
				return
			}
			defer handle.Close()

			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			closeAfter := time.After(time.Duration(runDuration) * time.Second)
			statics := make(map[string](*common.Metric))
			for packet := range packetSource.Packets() {
				select {
				case <-closeAfter:
					return
				default:
					m := printPacketInfo(packet, ifname, ifip)
					if m != nil {
						if existMetric, ok := statics[m.Key()]; !ok {
							statics[m.Key()] = m
							if len(statics) > maxPacketSize {
								for _, lm := range statics {
									mu.Lock()
									L = append(L, lm)
									mu.Unlock()
								}
								return
							}
						} else {
							if v, ok := existMetric.Value.(int); ok {
								existMetric.Value = v + 1
							}
						}
					}
				}
			}
		}(iface.Name, IPAddress)

	}
	wg.Wait()
	return
}

func printPacketInfo(packet gopacket.Packet, ifname string, ifip string) *common.Metric {
	// DNS layer recoder
	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer != nil {
		dns, ok := dnsLayer.(*layers.DNS)
		if ok {
			if (dns.ANCount == 0 && dns.ResponseCode > 0) || (dns.ANCount > 0) {
				for _, q := range dns.Questions {
					return toMetric("pcap.dns", 1, map[string]string{"Question": string(q.Name), "Type": q.Type.String(), "interface": ifname, "ip": ifip})
				}
			}
		}
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, ok := ipLayer.(*layers.IPv4)
		if ok {
			dstIP := ip.DstIP.String()
			srcIP := ip.SrcIP.String()
			return toMetric("pcap.ipv4", 1, map[string]string{"DstIP": dstIP, "SrcIP": srcIP, "interface": ifname, "ip": ifip})
		}
	}
	return nil
}
