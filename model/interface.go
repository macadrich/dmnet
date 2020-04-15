package model

import "net"

// P2PIFServer -
type P2PIFServer interface {
	Status()
	Stop()
	Listen()
	CreateConn(net.Conn, net.Addr) (Conn, error)
	OnMessage(func([]byte))
}

// IFServer tcp server
type IFServer interface {
	Addr() string
	Stop()
	Listen()
	CreateConn(net.Conn, net.Addr) (Conn, error)
	OnMessage(func([]byte))
}

// IFClient tcp client
type IFClient interface {
	GetServer() IFServer
	GetPeer() *Peer
	SetPeer(*Peer)
	GetSelf() *Peer
	Stop()
	OnMessage(func([]byte))
}
