package discovery

import (
	"fmt"
	"github.com/Mortimor1/mikromon-discovery/pkg/utils"
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
)

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func scan(subnets []string) ([]int64, error) {
	var alives []int64

	for _, subnet := range subnets {
		hosts, _ := Hosts(subnet)
		for _, host := range hosts {
			p := fastping.NewPinger()
			ra, err := net.ResolveIPAddr("ip4:icmp", host)
			if err != nil {
				return nil, err
			}
			p.AddIPAddr(ra)

			p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
				alives = append(alives, utils.IpToInt(addr.IP))
				fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
			}

			err = p.Run()
			if err != nil {
				return nil, err
			}
		}
	}
	return alives, nil
}
