package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/macadrich/dmnet/tcp"
)

func rndzServer() {
	log.Println("[ SERVER ]")
	s, err := tcp.New(tcp.DMNETSERVER, "0.0.0.0:9001")
	if err != nil {
		panic(err)
	}

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	log.Print(<-exit)
	s.Stop()
	log.Println("Done.")
}

func p2pClient() {
	log.Println("[ P2P ]")
	c, err := tcp.New(tcp.DMNETP2P, "175.176.79.16:44246")
	if err != nil {
		panic(err)
	}

	c.SignalInterupt()

}

func main() {
	p2pClient()
	//rndzServer()
}
