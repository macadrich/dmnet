package p2p

import (
	"testing"
)

func TestP2PClient(t *testing.T) {
	peer, err := NewP2P("adriel", "0.0.0.0:9001")
	if err != nil {
		t.Errorf("Unable to initialize p2p connection! %v", err)
	}

	err = peer.StartP2P()
	if err != nil {
		t.Errorf("Unable to establish p2p connection! %v", err)
	}

	peer.Stop()
}
