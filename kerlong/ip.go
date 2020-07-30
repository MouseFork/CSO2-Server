package kerlong

import (
	"errors"
	"fmt"
	"net"
)

//IPToUint32 把IP转换成4字节uint
func IPToUint32(s string) (uint32, error) {
	var ip uint32
	ipobj := net.ParseIP(s)
	if ipobj == nil {
		return ip, errors.New("Prase IP error !")
	} else {
		ip |= uint32(ipobj[12]) << 24
		ip |= uint32(ipobj[13]) << 16
		ip |= uint32(ipobj[14]) << 8
		ip |= uint32(ipobj[15])
	}
	return ip, nil
}

//SlideIP 切割IP，找到：位置
func SlideIP(s string) int {
	for k, v := range s {
		if v == ':' {
			return k
		}
	}
	return 0
}

//GetIP 获取IP
func GetIP() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return "Error"
}

//IsSameLan 判断两个IP是否处于同一局域网
func IsSameLan(a []byte, b []byte) bool {
	idx, i := 0, 0
	for {
		if len(a) <= i || len(b) <= i {
			return false
		}
		if a[i] == b[i] {
			if a[i] == '.' {
				idx++
				if idx == 3 {
					return true
				}
			}
		} else {
			return false
		}
		i++
	}
}
