package main

import (
	"dmnet/p2p"
	"dmnet/tcp"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func rndzServer() {
	s, err := tcp.New("server", "0.0.0.0:9001")
	if err != nil {
		panic(err)
	}
	s.Listen()

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	log.Print(<-exit)
	log.Println("Done.")
}

func p2pClient() {
	c, err := p2p.NewP2P("adriel", "0.0.0.0:9001")
	if err != nil {
		panic(err)
	}

	c.StartP2P()

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	log.Print(<-exit)
	log.Println("Done.")
}

func main() {
	p2pClient()
	//rndzServer()
}
