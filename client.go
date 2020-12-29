package main

import (
	"cp_cloud/lib"
	"flag"
	"github.com/prometheus/common/log"
	"net"
	"strings"
	"sync"
	"time"
)

//var lastCopiedString string

func main() {
	flag.Parse()

	// // my self define password
	//mySelfPassword := "m2y3_4p5a6s7s8w1o23rdea023_d13d1"
	//lib.ServerAuthFlag = &mySelfPassword

	if len(*lib.ServerAuthFlag) > 32 {
		log.Warn("The server auth key size cannot more than 32.")
		return
	}

	serverIpCh := make(chan string)
	go findServer(serverIpCh)
	serverIp, ok := <-serverIpCh
	if !ok {
		log.Error("Cannot find the server.")
		return
	}
	log.Info("The server ip is ", serverIp)

	lib.StartClient(serverIp)

}

func findServer(ch chan string) {
	ips, err := lib.GetClientIps()
	if err != nil {

	}
	var wg sync.WaitGroup

	for _, ipNet := range ips {
		wg.Add(1)
		multicastIp := lib.GetMultiCastAddr(ipNet)
		go multiCast(multicastIp, ch, &wg)
	}

	wg.Wait()
	close(ch)
}

func multiCast(multiCastIP string, ch chan string, wg *sync.WaitGroup) {
	ip := net.ParseIP(multiCastIP)
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: *lib.DiscoveryServiceFlag}
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		log.Error("[multiCast:1]", err)
		wg.Done()
		return
	}
	err = conn.SetDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		log.Error("[multiCast:2]", err)
		wg.Done()
		return
	}

	n, err := conn.WriteToUDP([]byte(*lib.ClientHelloFlag), dstAddr)
	if err != nil {
		log.Error("[multiCast:3]", err)
		wg.Done()
		return
	}
	data := make([]byte, 1024)

	n, addr, err := conn.ReadFrom(data)
	if err != nil {
		log.Error("[multiCast:4]", err)
		wg.Done()
		return
	}
	//fmt.Println("add => ", addr.String())
	//fmt.Print(addr.String())

	if string(data[:n]) == *lib.ServerHelloFlag {
		addrSpArr := strings.Split(addr.String(), ":")
		ch <- addrSpArr[0]
	}

	err = conn.Close()
	if err != nil {
		log.Error("[multiCast:5]", err)
	}
	wg.Done()
	return
}
