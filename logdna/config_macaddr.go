package logdna

import (
	"net"
	"strings"
)

// ref: https://socketloop.com/tutorials/golang-get-local-ip-and-mac-address
func getMacAndIP() (mac, ip string) {
	ipAddr := getIPAddr()
	if ipAddr == "" {
		return "", ""
	}
	hwName := getNetworkHardwareName(ipAddr)
	if hwName == "" {
		return "", ipAddr
	}

	netif, err := net.InterfaceByName(hwName)
	if err != nil {
		return "", ipAddr
	}

	return netif.HardwareAddr.String(), ipAddr
}

func getIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		ipnet, ok := address.(*net.IPNet)
		switch {
		case !ok,
			ipnet.IP.IsLoopback(),
			ipnet.IP.To4() == nil:
			continue
		}
		return ipnet.IP.String()
	}
	return ""
}

func getNetworkHardwareName(ipAddr string) string {
	netifs, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, netif := range netifs {
		addrs, err := netif.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if strings.Contains(addr.String(), ipAddr) {
				return netif.Name
			}
		}
	}
	return ""
}
