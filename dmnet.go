package dmnet

// RNDZServer rendezvous server
type RNDZServer interface {
	P2PEnable(bool)
	StartP2P() error
	Listen()
	Stop()
}

// Network network config
type Network struct {
	Addr string
	Mode string
}
