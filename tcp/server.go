package tcp

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/macadrich/dmnet/model"
	"github.com/macadrich/dmnet/util"
)

// P2PServer -
type P2PServer struct {
	sconn *net.TCPListener
	conns model.Conns
	wg    *sync.WaitGroup
	send  chan *model.Payload
	exit  chan bool
}

// Stop -
func (s *P2PServer) Stop() {
	close(s.exit)
	s.wg.Wait()
	log.Print("TCP Server exited")
}

// Status -
func (s *P2PServer) Status() {
	log.Println("IP:", s.sconn.Addr().String())
}

// SignalInterupt -
func (s *P2PServer) SignalInterupt() {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	log.Print(<-exit)
	s.Stop()
}

// p2psender -
func (s *P2PServer) p2psender() {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		select {
		case <-s.exit:
			log.Print("exiting TCP sender")
			return
		case p := <-s.send:
			if p != nil {
				log.Println("[p2psender()] Send:", string(p.Bytes), p.Addr.String())

				conn := s.conns[p.Addr.String()]
				c := conn.GetConn()

				c.Write(p.Bytes)
				log.Println("Send success:")
			}
		}
	}
}

// NewP2PServer -
func NewP2PServer(saddr *net.TCPAddr) (*P2PServer, error) {
	listener, err := net.ListenTCP("tcp", saddr)
	if err != nil {
		return nil, err
	}

	return &P2PServer{
		sconn: listener,
		conns: make(model.Conns),
		wg:    &sync.WaitGroup{},
		send:  make(chan *model.Payload, 100),
		exit:  make(chan bool),
	}, nil
}

// Listen -
func (s *P2PServer) Listen() {
	log.Println("Listening on", s.sconn.Addr())

	go s.p2psender()

	for {

		conn, err := s.sconn.Accept()
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
		log.Printf("New Connection: %v", tcpAddr)

		s.conns[tcpAddr.String()] = c
		s.wg.Add(1)
		go s.receive(conn)
	}
}

func (s *P2PServer) serve(b []byte, c net.Conn) {
	defer s.wg.Done()
	msg := &model.Message{}
	m, err := util.RecvMessage(msg, b)
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("Receive:", m, "From:", c.RemoteAddr().String())

	conn := s.conns[c.RemoteAddr().String()]
	conn.Send(&model.Message{
		Type:    "text",
		Content: "[serve] Hello: " + c.RemoteAddr().String(),
	})
}

func (s *P2PServer) receive(c net.Conn) {
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

// NewServer -
func NewServer(addr *net.TCPAddr) (*P2PServer, error) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &P2PServer{
		sconn: listener,
		conns: make(model.Conns),
		wg:    &sync.WaitGroup{},
		send:  make(chan *model.Payload, 100),
		exit:  make(chan bool),
	}, nil
}

// CreateConn -
func (s *P2PServer) CreateConn(conn net.Conn, sAddr net.Addr) (model.Conn, error) {
	if sAddr == nil {
		return nil, errors.New("Conns addr must not be nil")
	}

	tcpAddr, ok := sAddr.(*net.TCPAddr)
	if !ok {
		return nil, errors.New("could not assert net.Addr to *net.UDPAddr")
	}

	log.Println("Server address:", tcpAddr.String())

	c := model.NewPeerConn(conn, s.send, tcpAddr)
	s.conns[tcpAddr.String()] = c
	return c, nil
}
