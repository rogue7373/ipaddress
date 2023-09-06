package ipaddress

import (
	"errors"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

func GetIPAddress() ([]string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ipconfig")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("ifconfig")
	} else {
		return nil, errors.New("unsupported operating system")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	interfaces := strings.Split(string(output), "\n\n")
	var ipAddresses []string

	for _, iface := range interfaces {
		if len(iface) > 0 && !strings.Contains(iface, "loopback") && !strings.Contains(iface, "virtual network") {
			addrs, err := net.InterfaceAddrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						ip := strings.TrimSpace(ipNet.IP.String())
						if len(ip) > 0 {
							ipAddresses = append(ipAddresses, ip)
						}
					}
				}
			}
		}
	}

	return ipAddresses, nil
}
