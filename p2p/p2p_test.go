package p2p

import (
	"testing"

	"github.com/macadrich/dmnet/p2p"
)

func TestP2PClient(t *testing.T) {
	peer, err := p2p.NewP2P("adriel", "0.0.0.0:9001")
	if err != nil {
		t.Errorf("Unable to initialize p2p connection! %v", err)
	}

	err = peer.StartP2P()
	if err != nil {
		t.Errorf("Unable to establish p2p connection! %v", err)
	}
}
