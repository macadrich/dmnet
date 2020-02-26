package tcp

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"strconv"
)

// SendMessage -
func SendMessage(msg *Message) ([]byte, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// RecvMessage -
func RecvMessage(recv []byte) (*Message, error) {
	msg := &Message{}
	err := json.Unmarshal(recv, msg)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return msg, nil
}

// GenPort generate random port
func GenPort() string {
	return ":" + strconv.Itoa(rand.Intn(65535-10000)+10000)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func Do(n int, f func() error) error {
	var e error
	for i := 0; i < n; i += 1 {
		if e = f(); e == nil {
			return e
		}
	}
	return e
}
