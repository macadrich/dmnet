package tcp

import (
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

		return server, nil
	case "client":
		return nil, nil
	default:
		return nil, errors.New("unknown network mode")
	}
}
