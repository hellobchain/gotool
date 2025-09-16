package gnet

import (
	"net"
)

// IntranetIP 获取第一块内网 IPv4
func IntranetIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil {
				if ip[0] == 10 || ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 || ip[0] == 192 && ip[1] == 168 {
					return ip.String(), nil
				}
			}
		}
	}
	return "", nil
}

// IsPortOpen 检测端口是否开放
func IsPortOpen(addr string) bool {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
