package tcp

import (
	"dmnet/util"
	"errors"
	"net"

	"github.com/macadrich/dmnet"
)

// New network factory mode
func New(mode, address string) (dmnet.RNDZServer, error) {
	switch mode {
	case "server":
		tcpaddr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}

		server, err := NewTCPServer(tcpaddr, nil)
		if err != nil {
			return nil, err
		}

		server.P2PEnable(false)

		return server, nil
	case "client":
		return nil, nil
	case "p2p":
		saddr, err := net.ResolveTCPAddr("tcp", serveraddr)
		if err != nil {
			return nil, err
		}

		caddr, err := net.ResolveTCPAddr("tcp", util.GenPort())

		s, err := NewTCPServer(caddr, saddr)
		if err != nil {
			return nil, err
		}

		client, err := NewTCPClient(username, s)
		if err != nil {
			return nil, err
		}

		client.P2PEnable(true)
		return client, nil
	default:
		return nil, errors.New("unknown network mode")
	}
}
