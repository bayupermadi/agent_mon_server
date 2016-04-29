package main

import (
    "log"
    "errors"
	"fmt"
	"net"
	"os"
	"encoding/json"

    linuxproc "github.com/c9s/goprocinfo/linux"
)

func cpu_usage() uint64{
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
	    log.Fatal("stat read fail")
	}

	s := stat.CPUStatAll 
	a := s.User + s.Nice + s.System + s.IOWait + s.IRQ + s.SoftIRQ  + s.Steal + s.Guest + s.GuestNice + s.Idle
	c := s.User + s.Nice + s.System + s.IOWait + s.IRQ + s.SoftIRQ  + s.Steal + s.Guest + s.GuestNice

	d := c*100/a


   	//fmt.Println(d)
   	return d

}

func mem_usage() uint64{
	MemInfo, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
	    log.Fatal("meminfo read fail")
	}

	memtotal := MemInfo.MemTotal
	memactive := MemInfo.Active

	usage_util := memactive*100/memtotal

	//fmt.Println(usage_util)

	return usage_util
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
	//cpu_usage()
	//mem_usage()

	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("ip_address : ", ip)


    hostname, err := os.Hostname()

    if err != nil {
        panic(err)
    }

    cpu_usages := cpu_usage()
    mem_usages := mem_usage()

    //fmt.Println("hostname : ", hostname)

    type AMonS struct {
		Hoetname	string
		IP 			string	
		MemUsage 	uint64
		CPUUsage 	uint64
	}
	mons := AMonS{
		Hoetname: hostname,
		IP:   ip,
		MemUsage: mem_usages,
		CPUUsage: cpu_usages,
	}
	b, err := json.MarshalIndent(mons, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}