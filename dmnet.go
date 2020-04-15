package dmnet

// DMNet rendezvous server
type DMNet interface {
	Status()
	Stop()
	SignalInterupt()
	OnMessage(func([]byte))
}

// Network network config
type Network struct {
	Addr string
	Mode string
}
