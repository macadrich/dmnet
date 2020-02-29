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

// NewTCPClient initialize with username
// for client and server as p2p
func NewTCPClient(username, address string) (*Client, error) {

	nd, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	// set username peer
	self := &model.Peer{Username: username}

	// initialize peer
	p := &model.Peer{}

	return &Client{
		c:                  nd,
		self:               self,
		peer:               p,
		registeredCallback: func(model.IFClient) {},
		messageCallback:    func(model.IFClient, string) {},
	}, nil
}

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

// StartP2P start peer to peer connection
func (c *Client) StartP2P() error {
	listener, err := net.Listen("tcp", util.GenPort())
	if err != nil {
		return nil
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return nil
		}
		go handleRequest(conn, c.c)
	}

	/*
		s := c.GetServer()

		sConn, err := s.CreateConn(c.sAddr)
		if err != nil {
			return err
		}

		// set conn
		c.SetServerConn(sConn)

		// start listening
		go s.Listen()

		// send greeting message to server
		sConn.Send(&model.Message{
			Type:    "connect",
			Content: c.GetServer().Addr(),
		})
	*/

	return nil
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
