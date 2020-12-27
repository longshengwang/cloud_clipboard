package lib

import (
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
)

func IpAddrToInt(ipAddr string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ipAddr).To4())
	return ret.Int64()
}

func IntToIpAddr(ipInt int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ipInt>>24), byte(ipInt>>16), byte(ipInt>>8), byte(ipInt))
}

func GetMultiCastAddr(ipNet *net.IPNet) string {
	m, l := ipNet.Mask.Size()
	multi := strings.Repeat("1", l-m)

	multiInt, _ := strconv.ParseInt(multi, 2, 64)
	ipInt := IpAddrToInt(ipNet.IP.String())
	//println(l-m, multi, multiInt, ipInt)
	return IntToIpAddr(ipInt | multiInt)
}
