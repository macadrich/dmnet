package p2p

import (
	"net"

	"github.com/macadrich/dmnet/tcp"
)

// P2P represent client and server address
type P2P struct {
	*tcp.Client
	sAddr *net.TCPAddr
}

// NewP2P initialize peer to peer connection
func NewP2P(username string, serveraddr string) (*P2P, error) {
	saddr, err := net.ResolveTCPAddr("tcp", serveraddr)
	if err != nil {
		return nil, err
	}

	caddr, err := net.ResolveTCPAddr("tcp", tcp.GenPort())

	s, err := tcp.NewTCPServer(caddr, saddr)
	if err != nil {
		return nil, err
	}

	c, err := tcp.NewTCPClient(username, s)
	if err != nil {
		return nil, err
	}

	p2p := &P2P{
		Client: c,
		sAddr:  saddr,
	}

	return p2p, nil
}

// StartP2P start peer to peer connection
func (p2p *P2P) StartP2P() error {
	s := p2p.GetServer()

	sConn, err := s.CreateConn(p2p.sAddr)
	if err != nil {
		return err
	}

	// set conn
	p2p.SetServerConn(sConn)

	// start listening
	go s.Listen()

	// send greeting message to server
	sConn.Send(&tcp.Message{
		Type:    "connect",
		Content: p2p.GetServer().Addr(),
	})

	return nil
}
