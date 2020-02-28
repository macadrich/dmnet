package tcp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/macadrich/dmnet/model"
	"github.com/macadrich/dmnet/util"
)

// Server -
type Server struct {
	sAddr   *net.TCPAddr
	p2paddr string
	p2pc    net.Listener
	conns   model.Conns
	isP2P   bool
	wg      *sync.WaitGroup
	send    chan *model.Payload
	exit    chan bool
}

// NewTCPServer -
func NewTCPServer(mode string, addr *net.TCPAddr, saddr *net.TCPAddr) (*Server, error) {
	var m bool

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	if mode == DMNETP2P {
		m = true
	}

	port := listener.Addr().(*net.TCPAddr).Port
	paddr := util.GetLocalIP() + ":" + fmt.Sprintf("%d", port)

	return &Server{
		sAddr:   saddr,             // p2p server address to connect
		p2paddr: paddr,             // p2p address and port
		p2pc:    listener,          // p2p listener
		conns:   make(model.Conns), // p2p connections
		isP2P:   m,
		wg:      &sync.WaitGroup{},
		send:    make(chan *model.Payload, 100),
		exit:    make(chan bool),
	}, nil
}

func (s *Server) sender() {
	s.wg.Add(1)
	defer s.wg.Done()
	log.Println("sender()")
	if s.isP2P { // peer to peer only
		s.p2precv()
	} else {
		s.listenrecv()
	}
}

// Status -
func (s *Server) Status() {
	log.Println("IP:", s.sAddr.String())
	log.Println("Mode:", s.isP2P)
}

// listenrecv -
func (s *Server) listenrecv() {
	log.Println("listenrecv()")
	for {
		select {
		case <-s.exit:
			log.Print("exiting TCP sender")
			return
		case p := <-s.send:
			if p != nil {
				log.Println("Send:", string(p.Bytes))
				c := s.conns[p.Addr.String()]

				n, err := io.WriteString(c.GetTCPConn(), string(p.Bytes))
				if err != nil {
					return
				}
				log.Println("write:", n)

			}
		}
	}
}

// p2precv -
func (s *Server) p2precv() {
	conn, err := net.Dial("tcp", s.sAddr.String())
	if err != nil {
		log.Print("can't connect to server!")
		log.Printf("%v", err)
		return
	}
	defer conn.Close()

	for {
		select {
		case <-s.exit:
			log.Print("exiting TCP sender")
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
	log.Println("Listening on", s.Addr())
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

		c := model.NewPeerConn(conn, s.send, tcpAddr)
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

func (s *Server) serve(b []byte, c net.Conn) {
	defer s.wg.Done()
	msg := &model.Message{}
	m, err := util.RecvMessage(msg, b)
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("Receive:", m)
	conn := s.conns[c.RemoteAddr().String()]
	conn.Send(&model.Message{
		Type:    "text",
		Content: "Hello client!",
	})
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
		go s.serve(buf[:n], c)
	}
}

// CreateConn create connection to server
func (s *Server) CreateConn(sAddr net.Addr) (model.Conn, error) {
	if sAddr == nil {
		return nil, errors.New("Conns addr must not be nil")
	}

	tcpAddr, ok := sAddr.(*net.TCPAddr)
	if !ok {
		return nil, errors.New("could not assert net.Addr to *net.UDPAddr")
	}

	conn := s.conns[tcpAddr.String()]
	c := model.NewPeerConn(conn.GetTCPConn(), s.send, tcpAddr)
	s.conns[sAddr.String()] = c

	return c, nil
}
