package lib

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"github.com/atotto/clipboard"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lastCopiedString string

func StartClient(serviceIp string) {
	log.Println("Ready to start the client.")
	//ch := make(chan string)
	//go loopGetTextFromClipBoard(ch)
	conn, err := net.Dial("tcp", serviceIp+":"+strconv.Itoa(*ServerPortFlag))
	if err != nil {
		log.Fatalln(err)
		return
	}
	go handleServerConnection(conn)

	log.Println("Start the client successfully.")
	lastCopiedString, _ = clipboard.ReadAll()
	for true {
		n, _ := clipboard.ReadAll()
		if n != lastCopiedString {
			lastCopiedString = n
			_, err = conn.Write(GenConnByte(lastCopiedString))
			if err != nil {
				if err == io.EOF {
					log.Println("Remote Connect is closed, so end the client app.")
				} else {
					log.Println("Error:", err)
				}
				break
			}
		}
		time.Sleep(time.Duration(300) * time.Millisecond)
	}

	//
	//for s := range ch {
	//	_, err = conn.Write(GenConnByte(s))
	//	if err != nil {
	//		if err == io.EOF {
	//			log.Fatalln("Remote Connect is closed, so end the client app.")
	//			close(ch)
	//		} else {
	//			log.Println("Error:", err)
	//			close(ch)
	//		}
	//	}
	//}
}

func handleServerConnection(conn net.Conn) {
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
				log.Println("Conn is closed.")
			} else {
				log.Println("read conn with err: ", err)
			}
			break
		} else {
			buffer.Write(tmp[0:n])

			allLen := buffer.Len()
			allLenBack := allLen
			tmpBufferBytes := buffer.Bytes()
			scannerObj := bufio.NewScanner(buffer)
			scannerObj.Split(PacketSlitFunc)
			for scannerObj.Scan() {
				splitData := scannerObj.Bytes()
				allLen -= len(splitData)
				//fmt.Println("recv: ", string(splitData[8:]))
				encryptContent := string(splitData[8:])
				desContent := AesDecrypt(encryptContent, *ServerAuthFlag)
				lastCopiedString = desContent
				err := clipboard.WriteAll(desContent)
				if err != nil {
					log.Println("Copy to the clipboard with err:", err)
				}
			}
			if allLen > 0 {
				buffer.Write(tmpBufferBytes[allLenBack-allLen:])
			}
		}
	}
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
	ips, err := GetClientIps()
	if err != nil {

	}
	var wg sync.WaitGroup

	for _, ipNet := range ips {
		wg.Add(1)
		multicastIp := GetMultiCastAddr(ipNet)
		go multiCast(multicastIp, ch, &wg)
	}

	wg.Wait()
	close(ch)
}

func multiCast(multiCastIP string, ch chan string, wg *sync.WaitGroup) {
	ip := net.ParseIP(multiCastIP)
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: *DiscoveryServiceFlag}
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

	n, err := conn.WriteToUDP([]byte(*ClientHelloFlag), dstAddr)
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

	if string(data[:n]) == *ServerHelloFlag {
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

func cAuth(conn net.Conn) error {
	publicKeyByteEncrypt := make([]byte, 1024)
	n, e := conn.Read(publicKeyByteEncrypt)
	if e != nil {
		return e
	}
	var publicKeyP *rsa.PublicKey
	publicKeyP, e = x509.ParsePKCS1PublicKey(publicKeyByteEncrypt[:n])
	if e != nil {
		log.Println("error:", e)
		return e
	}

	passwdWithRandomKey := GenPasswordWithRandomKey(*ServerAuthFlag, 20)
	cipherBytes, e := rsa.EncryptPKCS1v15(rand.Reader, publicKeyP, []byte(passwdWithRandomKey))
	if e != nil {
		log.Println("error:", e)
		return e
	}
	n, e = conn.Write(cipherBytes)
	if e != nil {
		log.Println("error:", e)
		return e
	}
	return nil
}
