package dmnet

// P2P peer to peer client
type P2P interface {
	StartP2P()
}

// RNDZServer rendezvous server
type RNDZServer interface {
	Listen()
	Stop()
}

// Network network config
type Network struct {
	Addr string
	Mode string
}
