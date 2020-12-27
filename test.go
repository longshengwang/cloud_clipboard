package main

import (
	"bufio"
	"bytes"
	"cp_cloud/lib"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	//println(lib.GetSuitablePassword("12345678901234567890"))

	orig := "hello world dasdf asdf as dfas dfa sdf a sdf"
	key := "adasdfasd12312asdfa1231"
	fmt.Println("原文：", orig)

	encryptCode := lib.AesEncrypt(orig, key)
	fmt.Println("密文：", encryptCode)

	decryptCode := lib.AesDecrypt(encryptCode, key)
	fmt.Println("解密结果：", decryptCode)

	//g,_:=lib.GenPublicPrivateKey()
	//fmt.Println(g.PublicKey)
	////fmt.Println(g.PrivateKey)
	//d :=  x509.MarshalPKCS1PublicKey(&g.PublicKey)
	//println(len(d))
	//e,_ := x509.ParsePKCS1PublicKey(d)
	//fmt.Println(*e)
	//MarshalPKCS1PrivateKey(g.PrivateKey)

	//ch := make(chan string)
	//go func() {
	//	time.Sleep(1 * time.Second)
	//	//close(ch)
	//	ch <- "aaa"
	//}()
	//
	//d, ok := <- ch
	//println(d, ok)

	//var magicNum = 0x66dad
	//magicBin := make([] byte, 4)
	//binary.BigEndian.PutUint32(magicBin, 0x66dad)
	//fmt.Println(magicBin)
	//
	//fmt.Println(strings.Split("1.1.1.1:90", ":")[0])
	//b := lib.IpAddrToInt("1.1.1.1")
	//var a int64 = 0b11111111
	//println(b)
	//println(a|b)
	//fmt.Println(lib.IntToIpAddr(a|b))

	//res, _ := net.LookupHost("baidu.com")
	//println(len(res))
	//for i := range res {
	//	println(res[i])
	//}
	//
	//port, _:=net.LookupPort("tcp", "telnet")
	//println(port)
	//
	//

	//ch1 := make(chan string)
	//go loopGetTextFromClipBoard(ch1)
	//go startServer()
	//go startClient(ch1)
	//time.Sleep(time.Duration(300) * time.Second)
	//test()
	//testBuf()

	//println(strconv.ParseInt("122", 10, 4))
	//getClientIp()
}

//
//func getClientIp() (string, error) {
//
//	//interfaces, err := net.Interfaces()
//	//for _, inter := range interfaces {
//	//	println(inter.Name)
//	//	inter.Addrs()
//	//	addrs, _ := inter.Addrs()
//	//	for _, addr := range addrs {
//	//		println("  ", addr.String())
//	//	}
//	//
//	//}
//	println("==================")
//
//	addrs, err := net.InterfaceAddrs()
//
//	if err != nil {
//		return "", err
//	}
//
//	for _, address := range addrs {
//		// 检查ip地址判断是否回环地址
//		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//			if ipnet.IP.To4() != nil {
//				//fmt.Println(lib.GetMultiCastAddr(ipnet))
//				fmt.Println(ipnet.IP.String(), ipnet.Mask.String(), ipnet.IP.Mask(ipnet.Mask))
//			}
//
//		}
//	}
//
//	return "", nil
//}

func testBuf() {
	buff := bytes.NewBuffer(nil)
	a := make([]byte, 4)
	//b := make([]byte, 4)
	//c := make([]byte, 4)
	binary.BigEndian.PutUint32(a, uint32(100))
	buff.Write(a)
	println(buff.Len())
	var data int32
	binary.Read(buff, binary.BigEndian, &data)
	println(data)
	println(buff.Len())
}

func test() {
	//bytes.NewBuffer(nil)
	// An artificial input source.
	const input = "1234 5678 12345678 123456789 1234567891 12345678911  1234567901234567890"
	scanner := bufio.NewScanner(strings.NewReader(input))
	// Create a custom split function by wrapping the existing ScanWords function.
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanWords(data, atEOF)
		if err == nil && token != nil {
			_, err = strconv.ParseInt(string(token), 10, 32)
		}
		return
	}
	// Set the split function for the scanning operation.
	scanner.Split(split)
	// Validate the input
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s", err)
	}
}
