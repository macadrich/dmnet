package tcp

import (
	"bufio"
	"dmnet/util"
	"fmt"
	"log"
	"net"

	"github.com/macadrich/dmnet/model"
)

// Client base client
type Client struct {
	c                  net.Conn
	self               *model.Peer
	peer               *model.Peer
	sAddr              *net.TCPAddr
	sConn              model.Conn // server TCPConn
	pConn              model.Conn // peer TCPConn
	registeredCallback func(model.IFClient)
	messageCallback    func(model.IFClient, string)
}

// P2PClient -
type P2PClient struct {
	c     net.Conn
	s     model.P2PIFServer
	sConn model.Conn
	sAddr *net.TCPAddr
	self  *model.Peer
	peer  *model.Peer
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

// SetServerConn -
func (c *P2PClient) SetServerConn(conn model.Conn) {
	c.sConn = conn
}

// NewTCPClient -
func NewTCPClient(username, saddress string) (*P2PClient, error) {

	// connect to server
	c, err := net.Dial("tcp", saddress)
	if err != nil {
		return nil, err
	}

	saddr, err := net.ResolveTCPAddr("tcp", util.GenPort())
	if err != nil {
		return nil, err
	}

	// listen server
	s, err := NewP2PServer(saddr)
	if err != nil {
		return nil, err
	}

	// set username peer
	self := &model.Peer{Username: username}

	// initialize peer
	p := &model.Peer{}

	return &P2PClient{
		c:     c,
		s:     s,
		sAddr: saddr,
		self:  self,
		peer:  p,
	}, nil
}

func receive(c net.Conn) {
	defer c.Close()
	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			log.Print(err)
			return
		}
		log.Println("message:", string(buf[:n]))
	}
}

// StartP2P start peer to peer connection
func (c *P2PClient) StartP2P() error {

	s := c.GetServer()

	// start listening
	go s.Listen()

	msg := &model.Message{
		Type:    "connect",
		Content: c.sAddr.String(),
	}
	b, _ := util.SendMessage(msg)
	c.c.Write(b)

	receive(c.c)

	// send greeting message to server
	// sConn.Send(&model.Message{
	// 	Type:    "connect",
	// 	Content: "test",
	// })

	return nil
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - -

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()

	msg := bufio.NewScanner(src)
	for msg.Scan() {
		fmt.Println("Message:", msg.Text())
		dest.Write([]byte(msg.Text() + "\n"))
		src.Write([]byte(msg.Text() + "\n"))
	}
}

func handleRequest(conn, nd net.Conn) {
	fmt.Println("new client")

	go copyIO(conn, nd)
	go copyIO(nd, conn)
}

// Status -
func (c *Client) Status() {
	log.Println("IP:", c.self.Addr)
}

// OnRegistered -
func (c *Client) OnRegistered(f func(model.IFClient)) {
	c.registeredCallback = f
}

// OnMessage -
func (c *Client) OnMessage(f func(model.IFClient, string)) {
	c.messageCallback = f
}

// SetServerConn -
func (c *Client) SetServerConn(conn model.Conn) {
	c.sConn = conn
}

// GetServer -
// func (c *Client) GetServer() model.IFServer {
// 	return c.server
// }

// GetSelf -
func (c *Client) GetSelf() *model.Peer {
	return c.self
}

// GetPeer -
func (c *Client) GetPeer() *model.Peer {
	return c.peer
}

// SetPeer -
func (c *Client) SetPeer(p *model.Peer) {
	c.peer = p
}

// Stop -
func (c *Client) Stop() {
	//c.server.Stop()
}
