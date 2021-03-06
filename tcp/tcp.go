package tcp

import (
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

		server, err := NewServer(tcpaddr)
		if err != nil {
			return nil, err
		}

		server.Listen()

		return server, nil
	case DMNETCLIENT:
		return nil, nil
	case DMNETP2P:

		client, err := NewTCPClient("username", address)
		if err != nil {
			return nil, err
		}

		client.StartP2P()
		return client, nil
	default:
		return nil, errors.New("unknown network mode")
	}
}
