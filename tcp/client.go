package tcp

// NewTCPClient initialize with username
// for client and server as p2p
func NewTCPClient(username string, server IFServer) (*Client, error) {
	// set username peer
	self := &Peer{Username: username}

	// initialize peer
	p := &Peer{}

	return &Client{
		self:               self,
		peer:               p,
		server:             server,
		registeredCallback: func(IFClient) {},
		messageCallback:    func(IFClient, string) {},
	}, nil
}

// OnRegistered -
func (c *Client) OnRegistered(f func(IFClient)) {
	c.registeredCallback = f
}

// OnMessage -
func (c *Client) OnMessage(f func(IFClient, string)) {
	c.messageCallback = f
}

// SetServerConn -
func (c *Client) SetServerConn(conn Conn) {
	c.sConn = conn
}

// GetServer -
func (c *Client) GetServer() IFServer {
	return c.server
}

// GetSelf -
func (c *Client) GetSelf() *Peer {
	return c.self
}

// GetPeer -
func (c *Client) GetPeer() *Peer {
	return c.peer
}

// SetPeer -
func (c *Client) SetPeer(p *Peer) {
	c.peer = p
}

// Stop -
func (c *Client) Stop() {
	c.server.Stop()
}
