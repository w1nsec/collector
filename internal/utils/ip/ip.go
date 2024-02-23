package ip

import (
	"fmt"
	"net"
)

func GetIPv4() (ipv4 []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	a := make([]string, 0)

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			fmt.Println(err)
			continue
		}
		if ip4 := ip.To4(); ip4 != nil && !ip.IsLoopback() {
			a = append(a, ip4.String())
		}
	}

	if len(a) == 0 {
		return nil, fmt.Errorf("can't get non loopback address")
	}

	return a, nil
}
