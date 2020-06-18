package helper

import (
	"io/ioutil"
	"net"
	"net/http"
)

func GetMyPublicIp() (string, error) {
	// https://gist.github.com/ankanch/8c8ec5aaf374039504946e7e2b2cdf7f
	url := "https://api.ipify.org?format=text"
	// fmt.Println("Getting IP address from 'ipify'...")

	resp, err := http.Get(url)
	if err != nil {
		return "nil", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "nil", err
	}

	return string(ip), nil
}

// IsPublicIP check if an IP is in public range
// Example: IsPublicIP(net.ParseIP("192.168.66.6"))
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}
