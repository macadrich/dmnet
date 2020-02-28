package tcp

import (
	"net"

	"github.com/macadrich/dmnet/model"
)

// Client base client
type Client struct {
	server             model.IFServer
	self               *model.Peer
	peer               *model.Peer
	sAddr              *net.TCPAddr
	sConn              model.Conn // server TCPConn
	pConn              model.Conn // peer TCPConn
	isP2P              bool
	registeredCallback func(model.IFClient)
	messageCallback    func(model.IFClient, string)
}

// NewTCPClient initialize with username
// for client and server as p2p
func NewTCPClient(username string, server model.IFServer) (*Client, error) {
	// set username peer
	self := &model.Peer{Username: username}

	// initialize peer
	p := &model.Peer{}

	return &Client{
		self:               self,
		peer:               p,
		server:             server,
		registeredCallback: func(model.IFClient) {},
		messageCallback:    func(model.IFClient, string) {},
	}, nil
}

// StartP2P start peer to peer connection
func (c *Client) StartP2P() error {
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

	return nil
}

// Listen -
func (c *Client) Listen() {

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
func (c *Client) GetServer() model.IFServer {
	return c.server
}

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
	c.server.Stop()
}
