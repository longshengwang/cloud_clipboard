package lib

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
)

var magicNum uint32 = 0x66dad

func GenConnByte(content string) []byte {
	var magic = make([]byte, 4)
	binary.BigEndian.PutUint32(magic, magicNum)

	encryptContent := AesEncrypt(content, *ServerAuthFlag)

	lenNum := make([]byte, 4)
	binary.BigEndian.PutUint32(lenNum, uint32(len(encryptContent)))

	data := []byte(encryptContent)

	packetBuf := bytes.NewBuffer(magic)
	packetBuf.Write(lenNum)
	packetBuf.Write(data)
	return packetBuf.Bytes()
}

func PacketSlitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {

	// 检查 atEOF 参数 和 数据包头部的四个字节是否 为 0x123456(我们定义的协议的魔数)
	//println(">>>>>>>>>>>packetSlitFunc params", len(data))
	if !atEOF && len(data) > 8 && binary.BigEndian.Uint32(data[:4]) == magicNum {
		var l int32
		// 读出 数据包中 实际数据 的长度(大小为 0 ~ 2^16)
		binary.Read(bytes.NewReader(data[4:8]), binary.BigEndian, &l)
		pl := int(l) + 8
		//println(">>>>>sliit <<< ", pl, " data len:", len(data))
		if pl <= len(data) {
			return pl, data[:pl], nil
		} else {
			return 0, nil, nil
		}
	}
	return
}

func GetClientIps() ([]*net.IPNet, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return nil, err
	}
	var ips []*net.IPNet

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//return ipnet.IP.String(), nil
				ips = append(ips, ipnet)
			}
		}
	}

	return ips, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABC1D2E3F4G5H6I7J8K9L0M_=NOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenPasswordWithRandomKey(password string, randomSize int) string {
	return RandSeq(randomSize) + SplitKey + password
}
