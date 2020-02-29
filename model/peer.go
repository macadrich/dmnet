package model

import (
	"encoding/base64"
	"net"

	"github.com/macadrich/dmnet/util"
)

// Endpoint -
type Endpoint struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// Peer -
type Peer struct {
	ID         string       `json:"id,omitempty"`
	Username   string       `json:"username,omitempty"`
	Endpoint   Endpoint     `json:"endpoint,omitempty"`
	PublicKey  string       `json:"publicKey,omitempty"`
	PrivateKey [32]byte     `json:"-"`
	Addr       *net.UDPAddr `json:"-"`
}

// GetPublicKey -
func (p *Peer) GetPublicKey() ([32]byte, error) {
	var key [32]byte
	bs, err := base64.StdEncoding.DecodeString(p.PublicKey)
	if err != nil {
		return key, err
	}
	copy(key[:], bs)
	return key, nil
}

// SetPublicKey -
func (p *Peer) SetPublicKey(key [32]byte) {
	p.PublicKey = base64.StdEncoding.EncodeToString(key[:])
}

// PeerConn represents peer connection payload and address
type PeerConn struct {
	send chan *Payload
	addr *net.TCPAddr
	conn net.Conn
}

// Protocol connection tcp protocol
func (tcp *PeerConn) Protocol() string {
	return "TCP"
}

// GetAddr client tcp remote address
func (tcp *PeerConn) GetAddr() net.Addr {
	return tcp.addr
}

// Send message send to peer
func (tcp *PeerConn) Send(msg *Message) error {
	b, err := util.SendMessage(msg)
	if err != nil {
		return err
	}

	tcp.send <- &Payload{Bytes: b, Addr: tcp.addr}
	return nil
}

// GetTCPConn -
func (tcp *PeerConn) GetTCPConn() net.Conn {
	return tcp.conn
}

// NewPeerConn -
func NewPeerConn(c net.Conn, send chan *Payload, addr *net.TCPAddr) *PeerConn {
	return &PeerConn{
		send: send,
		addr: addr,
		conn: c,
	}
}
