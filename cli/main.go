package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/macadrich/dmnet/tcp"
)

// func rndzServer() {
// 	s, err := tcp.New("server", "0.0.0.0:9001")
// 	if err != nil {
// 		panic(err)
// 	}
// 	s.Listen()

// 	exit := make(chan os.Signal)
// 	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
// 	log.Print(<-exit)
// 	log.Println("Done.")
// }

func p2pClient() {
	c, err := tcp.New("p2p", "0.0.0.0:9001")
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
