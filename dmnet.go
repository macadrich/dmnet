package dmnet

// DMNet rendezvous server
type DMNet interface {
	Status()
	Stop()
}

// Network network config
type Network struct {
	Addr string
	Mode string
}
