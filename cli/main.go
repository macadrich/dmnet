package main

import (
	"flag"
	"log"
	"os"

	"github.com/macadrich/dmnet/tcp"
)

const (
	serverTCPPORT = ":9001"
	clientTCPPORT = ":44246"
)

var (
	ipaddr  = "0.0.0.0"
	tcpmode = "server"
	mode    = flag.String("mode", "server", "server by default")
	addr    = flag.String("addr", "", "IP address of rendezvous server")
)

// setup rendezvous server connection
// @addr - rendezvous address listening
func rndzServer(addr string) {
	log.Println("[ SERVER ]")
	s, err := tcp.New(tcp.DMNETSERVER, addr)
	if err != nil {
		panic(err)
	}

	s.SignalInterupt()
}

// connect to peer
// @addr - address of peer client to connect
func p2pClient(addr string) {
	log.Println("[ P2P ]")
	c, err := tcp.New(tcp.DMNETP2P, addr)
	if err != nil {
		panic(err)
	}

	c.SignalInterupt()
}

func main() {
	flag.Parse()

	if len(os.Args) <= 2 || len(os.Args) > 3 {
		log.Printf("usage: -mode=server -addr=0.0.0.0 \n")
		return
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("MODE =", *mode)
	switch *mode {
	case "client":
		if *addr != "" {
			ipaddr = *addr + clientTCPPORT
			p2pClient(ipaddr)
		}
		return
	case "server":
		if *addr != "" {
			ipaddr = *addr + serverTCPPORT
			rndzServer(ipaddr)
		}
		return
	default:
		log.Printf("Unable to execute mode \"%s\" \n", *mode)
	}
}
