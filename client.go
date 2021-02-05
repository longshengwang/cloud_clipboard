package main

import (
	"cp_cloud/lib"
	"flag"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

//var lastCopiedString string

func main() {
	flag.Parse()

	loopInterval := 5
	loopCount := 0

	if len(*lib.ServerAuthFlag) > 32 {
		log.Fatalln("The server auth key size cannot more than 32.")
		return
	}
	for {
		serverIpCh := make(chan string)
		go findServer(serverIpCh)
		serverIp, ok := <-serverIpCh
		if !ok {
			loopCount++
			loopInterval = getLoopInterval(loopCount)
			log.Println("Cannot find the server. Go to sleep ", loopInterval, "s")
			time.Sleep(time.Duration(loopInterval) * time.Second)
			continue
		}
		log.Println("The server ip is ", serverIp)
		loopCount = 0
		loopInterval = 5
		lib.StartClient(serverIp)
		time.Sleep(time.Duration(loopInterval) * time.Second)
		log.Println("Oops! The server ", serverIp, " has been broken.")
	}
}

func getLoopInterval(count int) int {
	if count < 20 {
		return 5
	} else if count >= 20 && count < 100 {
		return 5 + count/2
	} else {
		return 60
	}
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
		log.Println("[multiCast:1]", err)
		wg.Done()
		return
	}
	err = conn.SetDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		log.Println("[multiCast:2]", err)
		wg.Done()
		return
	}

	n, err := conn.WriteToUDP([]byte(*lib.ClientHelloFlag), dstAddr)
	if err != nil {
		log.Println("[multiCast:3]", err)
		wg.Done()
		return
	}
	data := make([]byte, 1024)

	n, addr, err := conn.ReadFrom(data)
	if err != nil {
		log.Println("[multiCast:4]", err)
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
		log.Println("[multiCast:5]", err)
	}
	wg.Done()
	return
}
