package dmnet

// RNDZServer rendezvous server
type RNDZServer interface {
	P2PEnable(bool)
	StartP2P()
	Listen()
	Stop()
}

// Network network config
type Network struct {
	Addr string
	Mode string
}
