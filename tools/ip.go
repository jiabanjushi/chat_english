package tools

import (
	"errors"
	"github.com/ipipdotnet/ipdb-go"
	"net"
	"strings"
)

func ParseIp(myip string) *ipdb.CityInfo {
	db, err := ipdb.NewCity("./config/city.free.ipdb")
	if err != nil {
		return nil
	}
	db.Reload("./config/city.free.ipdb")
	c, err := db.FindInfo(myip, "CN")
	if err != nil {
		return nil
	}
	return c
}
func ParseIpNew(myip string) *CityInfo {
	var cityInfo = &CityInfo{}
	var ipstr string
	var infos []string
	p, _ := NewIpdb("./config/qqzeng-ip-utf8.dat")
	ipstr = p.Get(myip)
	infos = strings.Split(ipstr, "|")
	if infos[1] == "保留" {
		cityInfo.CountryName = "Internal network"
		return cityInfo
	}
	cityInfo.CountryName = infos[1]
	cityInfo.RegionName = infos[2]
	cityInfo.CityName = infos[3]
	return cityInfo
}
func GetServerIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

//获取出站IP地址
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
