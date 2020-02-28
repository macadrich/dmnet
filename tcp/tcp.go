package tcp

import (
	"dmnet/util"
	"errors"
	"net"

	"github.com/macadrich/dmnet"
)

// New network factory mode
func New(mode, address string) (dmnet.DMNet, error) {
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

		server.Listen()

		return server, nil
	case "client":
		return nil, nil
	case "p2p":
		saddr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}

		caddr, err := net.ResolveTCPAddr("tcp", util.GenPort())

		s, err := NewTCPServer(caddr, saddr)
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
