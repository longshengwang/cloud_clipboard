package main

import (
	"bufio"
	"bytes"
	"cp_cloud/lib"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"github.com/prometheus/common/log"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

type InnerMessage struct {
	connRemoteString string
	copiedContent    string
}

var globalConnMap sync.Map
var globalRsaKey *lib.RsaKey

func main() {
	flag.Parse()

	var err error
	globalRsaKey, err = lib.GenPublicPrivateKey()
	if err != nil {
		log.Error("Gen Rsa private/public key with err:", err)
		return
	}
	startServer()
}

func startServer() {
	shareCh := make(chan InnerMessage)
	port := *lib.ServerPortFlag
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Error("Cannot start the server at port ", port)
		return
	}

	go startDiscoveryService()

	go clientsContentShared(shareCh)

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn, shareCh)
	}
}

func clientsContentShared(ch chan InnerMessage) {
	for c := range ch {
		log.Info(c.connRemoteString, " [COPIED] ", c.copiedContent)
		globalConnMap.Range(func(key, value interface{}) bool {
			if key != c.connRemoteString {
				p, ok := value.(net.Conn)
				if ok {
					_, err := p.Write(lib.GenConnByte(c.copiedContent))
					if err != nil {
						log.Error(err)
					}
				}
			}
			return true
		})
	}
}

func handleConnection(conn net.Conn, shareCh chan InnerMessage) {
	log.Debug(conn.RemoteAddr().String(), " is Connect")
	//conn.Write()
	isClientAccess, err := auth(conn)
	if err != nil {
		conn.Close()
		return
	}
	if !isClientAccess {
		conn.Close()
		return
	}
	log.Info(conn.RemoteAddr().String(), " is connect with correct password.")

	globalConnMap.Store(conn.RemoteAddr().String(), conn)
	tmp := make([]byte, 1024)
	buffer := bytes.NewBuffer(nil)
	for {
		n, err := conn.Read(tmp[0:])
		buffer.Write(tmp[0:n])
		if err != nil {
			if err != io.EOF {
				log.Error("Read conn with err: ", err)

			} else {
				log.Error("Connection ", conn.RemoteAddr().String(), " is closed.")
			}
			connClosed(conn)
			break
		} else {
			allLen := buffer.Len()
			allLenBack := allLen
			tmpBufferBytes := buffer.Bytes()
			scannerObj := bufio.NewScanner(buffer)
			scannerObj.Split(lib.PacketSlitFunc)
			for scannerObj.Scan() {
				splitData := scannerObj.Bytes()
				allLen -= len(splitData)

				encryptContent := string(splitData[8:])
				desContent := lib.AesDecrypt(encryptContent, *lib.ServerAuthFlag)

				shareCh <- InnerMessage{conn.RemoteAddr().String(), desContent}
			}
			if allLen > 0 {
				buffer.Write(tmpBufferBytes[allLenBack-allLen:])
			}
		}
	}
}

func connClosed(conn net.Conn) {
	globalConnMap.Delete(conn.RemoteAddr().String())
}

func startDiscoveryService() {
	pc, err := net.ListenPacket("udp", ":" + strconv.Itoa(*lib.DiscoveryServiceFlag))
	if err != nil {

	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go serve(pc, addr, buf[:n])
	}
}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	//log.Info("Recv From Client:", string(buf))
	if string(buf) == *lib.ClientHelloFlag {
		pc.WriteTo([]byte(*lib.ServerHelloFlag), addr)
	}
}

func auth(conn net.Conn) (bool, error) {
	var len int
	var err error

	// send the public to client
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&globalRsaKey.PublicKey)
	len, err = conn.Write(publicKeyBytes)

	// read the encrypt key of (randomkey + password)
	authKeyBytesEncrypt := make([]byte, 1024)
	len, err = conn.Read(authKeyBytesEncrypt)
	if err != nil {
		log.Error("Cannot read auth key from client:", conn.RemoteAddr().String(), ". Error: ", err)
		return false, err
	}
	var authKeyBytes []byte
	authKeyBytes, err = rsa.DecryptPKCS1v15(rand.Reader, &globalRsaKey.PrivateKey, authKeyBytesEncrypt[:len])
	randomKeyAndPassword := string(authKeyBytes)
	splitIndex := strings.Index(randomKeyAndPassword, "]")
	if splitIndex == -1 {
		log.Error("Cannot find the split key from the key. Client:", conn.RemoteAddr().String())
		return false, nil
	}
	password := randomKeyAndPassword[splitIndex+1:]
	if password == *lib.ServerAuthFlag {
		return true, nil
	}
	log.Error("Password is not same with the server password. Client:", conn.RemoteAddr().String())
	return false, nil
}
