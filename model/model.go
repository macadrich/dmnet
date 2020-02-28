package model

import (
	"net"
)

// Message represents model for exchange message between peer
type Message struct {
	Type    string      `json:"type"`
	PeerID  string      `json:"peerID,omitempty"`
	Error   string      `json:"error,omitempty"`
	Content interface{} `json:"data,omitempty"`
	Encrypt bool        `json:"-"`
	addr    *net.TCPAddr
}

// Payload represent network payload
type Payload struct {
	Bytes []byte
	Addr  *net.TCPAddr
}

// Conn peer connection interface
type Conn interface {
	Send(*Message) error
	Protocol() string
	GetAddr() net.Addr
	GetTCPConn() net.TCPConn
}

// Conns map peer connection
type Conns map[string]Conn
