package tcp

import "net"

// IFServer tcp server
type IFServer interface {
	Addr() string
	Stop()
	Listen()
	CreateConn(net.Addr) (Conn, error)
}

// IFClient tcp client
type IFClient interface {
	GetServer() IFServer
	GetPeer() *Peer
	SetPeer(*Peer)
	GetSelf() *Peer
	Stop()
	OnRegistered(func(IFClient))
	OnMessage(func(IFClient, string))
}
