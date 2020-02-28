package tcp

import (
	"dmnet/util"
	"errors"
	"net"

	"github.com/macadrich/dmnet"
)

const (
	// DMNETSERVER as server
	DMNETSERVER = "server"
	// DMNETP2P as peer to peer
	DMNETP2P = "p2p"
	// DMNETCLIENT as client
	DMNETCLIENT = "client"
)

// New network factory mode
func New(mode, address string) (dmnet.DMNet, error) {
	switch mode {
	case DMNETSERVER:
		tcpaddr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}

		server, err := NewTCPServer(mode, tcpaddr, nil)
		if err != nil {
			return nil, err
		}

		server.Listen()

		return server, nil
	case DMNETCLIENT:
		return nil, nil
	case DMNETP2P:
		saddr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}

		caddr, err := net.ResolveTCPAddr("tcp", util.GenPort())

		s, err := NewTCPServer(mode, caddr, saddr)
		if err != nil {
			return nil, err
		}

		client, err := NewTCPClient("username", s)
		if err != nil {
			return nil, err
		}

		client.StartP2P()
		return client, nil
	default:
		return nil, errors.New("unknown network mode")
	}
}
