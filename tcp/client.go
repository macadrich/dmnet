package tcp

import (
	"dmnet/util"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/macadrich/dmnet/model"
)

// Client base client
type Client struct {
	c               net.Conn
	self            *model.Peer
	peer            *model.Peer
	sAddr           *net.TCPAddr
	sConn           model.Conn // server TCPConn
	pConn           model.Conn // peer TCPConn
	messageCallback func([]byte)
}

// P2PClient -
type P2PClient struct {
	c               net.Conn
	s               model.P2PIFServer
	sConn           model.Conn
	sAddr           *net.TCPAddr
	self            *model.Peer
	peer            *model.Peer
	messageCallback func([]byte)
}

// OnMessage -
func (c *P2PClient) OnMessage(f func([]byte)) {
	c.messageCallback = f
}

// GetServer -
func (c *P2PClient) GetServer() model.P2PIFServer {
	return c.s
}

// Status -
func (c *P2PClient) Status() {
	log.Println("IP:", c.self.Addr)
}

// Stop -
func (c *P2PClient) Stop() {
	c.s.Stop()
}

// SignalInterupt -
func (c *P2PClient) SignalInterupt() {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	log.Print(<-exit)
	c.Stop()
}

// SetServerConn -
func (c *P2PClient) SetServerConn(conn model.Conn) {
	c.sConn = conn
}

// NewTCPClient -
func NewTCPClient(username, saddress string) (*P2PClient, error) {
	var s *P2PServer
	serverAddr, err := net.ResolveTCPAddr("tcp", saddress)
	if err != nil {
		return nil, err
	}

	// connect to server
	c, err := net.Dial("tcp", saddress)
	if err != nil {
		return nil, err
	}

	err = util.GetPortConn(5, func() error {
		saddr, err := net.ResolveTCPAddr("tcp", util.GenPort())
		if err != nil {
			return err
		}

		// listen server
		s, err = NewP2PServer(saddr)
		if err != nil {
			return err
		}
		return err
	})

	if err != nil {
		return nil, err
	}

	// set username peer
	self := &model.Peer{Username: username}

	// initialize peer
	p := &model.Peer{}

	return &P2PClient{
		c:               c,
		s:               s,
		sAddr:           serverAddr,
		self:            self,
		peer:            p,
		messageCallback: func([]byte) {},
	}, nil
}

func (c *P2PClient) receive() {
	defer c.c.Close()
	for {
		buf := make([]byte, 1024)
		n, err := c.c.Read(buf)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			log.Print(err)
			return
		}
		c.messageCallback(buf[:n])
	}
}

// StartP2P start peer to peer connection
func (c *P2PClient) StartP2P() error {

	s := c.GetServer()

	// start listening
	go s.Listen()

	conn, err := s.CreateConn(c.c, c.sAddr)
	if err != nil {
		panic(err)
	}

	c.SetServerConn(conn)

	go c.receive()

	conn.Send(&model.Message{
		Type:    "connect",
		Content: c.c.RemoteAddr().String(),
	})

	return nil
}
