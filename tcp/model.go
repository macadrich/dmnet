package tcp

import (
	"net"
	"sync"
)

// Server represents tcp server handle connection
type Server struct {
	sAddr   *net.TCPAddr
	p2paddr string
	p2pc    net.Listener
	conns   Conns
	wg      *sync.WaitGroup
	send    chan *Payload
	exit    chan bool
}

// Client base client
type Client struct {
	server             IFServer
	self               *Peer
	peer               *Peer
	sAddr              *net.TCPAddr
	sConn              Conn // server TCPConn
	pConn              Conn // peer TCPConn
	registeredCallback func(IFClient)
	messageCallback    func(IFClient, string)
}

// Message represents model for exchange message between peer
type Message struct {
	Type    string      `json:"type"`
	PeerID  string      `json:"peerID,omitempty"`
	Error   string      `json:"error,omitempty"`
	Content interface{} `json:"data,omitempty"`
	Encrypt bool        `json:"-"`
	addr    *net.TCPAddr
}

// Conn peer connection interface
type Conn interface {
	Send(*Message) error
	Protocol() string
	GetAddr() net.Addr
}

// Conns map peer connection
type Conns map[string]Conn

// Payload represent network payload
type Payload struct {
	Bytes []byte
	Addr  *net.TCPAddr
}
