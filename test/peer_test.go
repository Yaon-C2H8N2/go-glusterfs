package test

import (
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/peer"
	"testing"
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
