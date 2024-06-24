package test

import (
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/eventListener"
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/peer"
	"testing"
	"time"
)

func TestProbing(t *testing.T) {
	hostnameList := []string{"node2", "node3"}
	for _, hostname := range hostnameList {
		err := peer.PeerProbe(hostname)
		if err != nil {
			t.Fatalf("Failed to probe %s: %s", hostname, err)
		}
	}
}

func TestDetachPeer(t *testing.T) {
	peers, err := peer.ListPeers()
	if err != nil {
		t.Fatalf("Failed to list peers: %s", err)
	}
	if len(peers) == 0 {
		t.Fatalf("No peers found. At least one peer should be present.")
	}
	for _, p := range peers {
		if p.Hostname != "node1" {
			err = p.Detach()
			if err != nil {
				t.Fatalf("Failed to detach peer %s: %s", p.Hostname, err)
			}
		}
	}
}

func TestPeerEvents(t *testing.T) {
	peerProbed := false
	peerRemoved := false
	//peerDisconnected := false
	//peerConnected := false

	peerUpdateHandler := func(event eventListener.PeerEvent) {
		if event.Type == eventListener.PEER_PROBED {
			peerProbed = true
		}
		if event.Type == eventListener.PEER_REMOVED {
			peerRemoved = true
		}
		//if event.Type == eventListener.PEER_DISCONNECTED {
		//	peerDisconnected = true
		//}
		//if event.Type == eventListener.PEER_CONNECTED {
		//	peerConnected = true
		//}
	}

	listener := eventListener.Default()
	listener.SetPollingTimeout(100)
	listener.OnPeerUpdate = peerUpdateHandler
	err := listener.Start()
	if err != nil {
		t.Fatalf("Failed to start event listener: %s", err)
	}
	// Wait for listener to start
	time.Sleep(200 * time.Millisecond)

	TestProbing(t)
	// Wait for volume create event
	time.Sleep(200 * time.Millisecond)
	if !peerProbed {
		t.Fatalf("Peer probe event not received")
	}

	t.Logf("Detaching peers")
	TestDetachPeer(t)
	// Wait for volume deletion event
	time.Sleep(200 * time.Millisecond)
	if !peerRemoved {
		t.Fatalf("Peer remove event not received")
	}
	err = listener.Stop()
	if err != nil {
		t.Fatalf("Failed to stop event listener: %s", err)
	}
}
