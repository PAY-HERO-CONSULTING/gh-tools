package address

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
)

func GetRequestMacAddress() (string, string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Infof("error retrieving network interfaces: [%+v]", err)
		return "", "", err
	}

	serverIP, err := getServerIP()
	if err != nil {
		log.Println("Error retrieving server IP:", err)
		return "", "", err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			logger.Infof("Error retrieving interface addresses: [%+v]", err)
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			ipAddress := fmt.Sprintf("%v", ipNet.IP)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				if ipNet.IP.Equal(serverIP) {
					if i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil) {
						// Skip locally administered addresses
						if i.HardwareAddr[0]&2 == 2 {
							continue
						}

						mac := formatMAC(i.HardwareAddr)
						return mac, ipAddress, nil
					}
				} else {
					logger.Infof("Connection from different network detected. IP: [%+v]: MAC Address: [%v]", ipNet.IP, formatMAC(i.HardwareAddr)) // can choose to send messages to admin
					// white list IP list
					return ipNet.IP.String(), formatMAC(i.HardwareAddr), nil
				}
			}
		}
	}

	return "", "", fmt.Errorf("no valid MAC address found") // redirect to login
}

func formatMAC(macAddr net.HardwareAddr) string {
	macAddress := ""
	for i, b := range macAddr {
		macAddress += fmt.Sprintf("%02x", b)
		if i < len(macAddr)-1 {
			macAddress += ":"
		}
	}
	return macAddress
}

func getServerIP() (net.IP, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4, nil
		}
	}

	return nil, fmt.Errorf("no IPv4 address found for the server")
}
