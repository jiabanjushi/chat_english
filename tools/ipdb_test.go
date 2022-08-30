package tools

import (
	"log"
	"strings"
	"testing"
)

func TestNewIpdb(t *testing.T) {
	var ip, ipstr string
	var infos []string
	p, _ := NewIpdb("../config/qqzeng-ip-utf8.dat")
	ip = "113.104.209.240"
	ipstr = p.Get(ip)
	infos = strings.Split(ipstr, "|")
	log.Println(infos)

	ip = "39.155.215.54"
	ipstr = p.Get(ip)
	infos = strings.Split(ipstr, "|")
	log.Println(infos)

	ip = "127.0.0.1"
	ipstr = p.Get(ip)
	infos = strings.Split(ipstr, "|")
	log.Println(infos)

	ip = "192.168.1.1"
	ipstr = p.Get(ip)
	infos = strings.Split(ipstr, "|")
	for k, v := range infos {
		log.Println(k, v)
	}
}
