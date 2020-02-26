package tcp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// NewTCPServer -
func NewTCPServer(addr *net.TCPAddr, saddr *net.TCPAddr) (*Server, error) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	paddr := getLocalIP() + ":" + fmt.Sprintf("%d", port)

	return &Server{
		sAddr:   saddr,       // p2p server address to connect
		p2paddr: paddr,       // p2p address and port
		p2pc:    listener,    // p2p listener
		conns:   make(Conns), // p2p connections
		wg:      &sync.WaitGroup{},
		send:    make(chan *Payload, 100),
		exit:    make(chan bool),
	}, nil
}

func (s *Server) sender() {
	s.wg.Add(1)
	defer s.wg.Done()

	conn, err := net.Dial("tcp", s.sAddr.String())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		select {
		case <-s.exit:
			log.Print("exiting UDP sender")
			return
		case p := <-s.send:
			if p != nil {
				log.Println("Send:", string(p.Bytes))
				conn.Write(p.Bytes)
			}
		}
	}
}

// Addr -
func (s *Server) Addr() string {
	return s.p2paddr
}

// Listen -
func (s *Server) Listen() {
	log.Println("listening...")
	go s.sender()
	for {

		conn, err := s.p2pc.Accept()
		if err != nil {
			log.Print(err) // print error and continue
			continue
		}

		tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
		if !ok {
			log.Print("could not assert net.Addr to *net.TCPAddr")
			return
		}

		c := NewPeerConn(s.send, tcpAddr)
		log.Printf("New Connection: %v", c.GetAddr())
		s.conns[tcpAddr.String()] = c

		s.receive(conn)
	}

}

// Stop -
func (s *Server) Stop() {
	close(s.exit)
	s.wg.Wait()
	log.Print("TCP Server exited")
}

func (s *Server) serve(b []byte) {
	defer s.wg.Done()
	m, err := RecvMessage(b)
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("Receive:", m)
}

func (s *Server) receive(c net.Conn) {
	defer c.Close()
	s.wg.Add(1)
	defer s.wg.Done()

	log.Println("Client:", c.RemoteAddr().String())
	for {
		select {
		case <-s.exit:
			log.Println("TCP exit [receive]")
			return
		default:
		}

		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			delete(s.conns, c.RemoteAddr().String())
			log.Print(err)
			return
		}

		s.wg.Add(1)
		go s.serve(buf[:n])
	}
}

// CreateConn create connection to server
func (s *Server) CreateConn(sAddr net.Addr) (Conn, error) {
	if sAddr == nil {
		return nil, errors.New("Conns addr must not be nil")
	}

	tcpAddr, ok := sAddr.(*net.TCPAddr)
	if !ok {
		return nil, errors.New("could not assert net.Addr to *net.UDPAddr")
	}

	c := NewPeerConn(s.send, tcpAddr)
	s.conns[sAddr.String()] = c

	return c, nil
}
