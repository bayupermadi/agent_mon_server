package main

import (
    "log"
    "errors"
	"fmt"
	"net"

    linuxproc "github.com/c9s/goprocinfo/linux"
)

func cpu_usage(){
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
	    log.Fatal("stat read fail")
	}

	s := stat.CPUStatAll 
	a := s.User + s.Nice + s.System + s.IOWait + s.IRQ + s.SoftIRQ  + s.Steal + s.Guest + s.GuestNice + s.Idle
	c := s.User + s.Nice + s.System + s.IOWait + s.IRQ + s.SoftIRQ  + s.Steal + s.Guest + s.GuestNice

	d := c*100/a


   	fmt.Println("cpu_usage : ",d,"%")

}

func mem_usage(){
	MemInfo, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
	    log.Fatal("meminfo read fail")
	}

	memtotal := MemInfo.MemTotal
	memactive := MemInfo.Active

	usage_util := memactive*100/memtotal

	fmt.Println("mem_usage : ",usage_util,"%")
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
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
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func main() {
	cpu_usage()
	mem_usage()

	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ip_address : ", ip)
}