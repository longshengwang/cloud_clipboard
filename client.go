package main

import (
	"bufio"
	"bytes"
	"cp_cloud/lib"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"github.com/atotto/clipboard"
	"github.com/prometheus/common/log"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lastCopiedString string

func main() {
	flag.Parse()

	serverIpCh := make(chan string)
	go findServer(serverIpCh)
	serverIp, ok := <-serverIpCh
	if !ok {
		log.Error("Cannot find the server.")
		return
	}
	log.Info("The server ip is ", serverIp)

	ch := make(chan string)
	go loopGetTextFromClipBoard(ch)
	startClient(ch, serverIp)

}

func startClient(ch chan string, serviceIp string) {
	conn, err := net.Dial("tcp", serviceIp+":"+strconv.Itoa(*lib.ServerPortFlag))
	if err != nil {
		log.Error(err)
		return
	}

	go handleServerConnection(conn, ch)
	for s := range ch {
		_, err = conn.Write(lib.GenConnByte(s))
		if err != nil {
			if err == io.EOF {
				log.Error("Remote Connect is closed, so end the client app.")
				close(ch)
			} else {
				log.Error("[54]error:", err)
				close(ch)
			}
		}
	}
}

func handleServerConnection(conn net.Conn, ch chan string) {
	err := cAuth(conn)
	if err != nil {
		conn.Close()
		return
	}
	tmp := make([]byte, 1024)
	buffer := bytes.NewBuffer(nil)
	for {
		n, err := conn.Read(tmp[0:])
		if err != nil {
			if err == io.EOF {
				log.Error("Conn is closed.")
			} else {
				log.Error("read conn with err: ", err)
			}
			break
		} else {
			buffer.Write(tmp[0:n])

			allLen := buffer.Len()
			allLenBack := allLen
			tmpBufferBytes := buffer.Bytes()
			scannerObj := bufio.NewScanner(buffer)
			scannerObj.Split(lib.PacketSlitFunc)
			for scannerObj.Scan() {
				splitData := scannerObj.Bytes()
				allLen -= len(splitData)
				//fmt.Println("recv: ", string(splitData[8:]))
				encryptContent := string(splitData[8:])
				desContent := lib.AesDecrypt(encryptContent, *lib.ServerAuthFlag)
				lastCopiedString = desContent
				err := clipboard.WriteAll(desContent)
				if err != nil {
					log.Error("Copy to the clipboard with err:", err)
				}
			}
			if allLen > 0 {
				buffer.Write(tmpBufferBytes[allLenBack-allLen:])
			}
		}
	}

	close(ch)
}

func loopGetTextFromClipBoard(ch chan string) {
	lastCopiedString, _ = clipboard.ReadAll()
	for true {
		n, _ := clipboard.ReadAll()
		if n != lastCopiedString {
			lastCopiedString = n
			ch <- lastCopiedString
		}
		time.Sleep(time.Duration(300) * time.Millisecond)
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

func cAuth(conn net.Conn) error {
	publicKeyByteEncrypt := make([]byte, 1024)
	n, e := conn.Read(publicKeyByteEncrypt)
	if e != nil {
		return e
	}
	var publicKeyP *rsa.PublicKey
	publicKeyP, e = x509.ParsePKCS1PublicKey(publicKeyByteEncrypt[:n])
	if e != nil {
		log.Error("error:", e)
		return e
	}

	passwdWithRandomKey := lib.GenPasswordWithRandomKey(*lib.ServerAuthFlag, 20)
	cipherBytes, e := rsa.EncryptPKCS1v15(rand.Reader, publicKeyP, []byte(passwdWithRandomKey))
	if e != nil {
		log.Error("error:", e)
		return e
	}
	n, e = conn.Write(cipherBytes)
	if e != nil {
		log.Error("error:", e)
		return e
	}
	return nil
}
