package eventListener

import (
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/peer"
)

const (
	PEER_PROBED       = "probed"
	PEER_REMOVED      = "removed"
	PEER_DISCONNECTED = "disconnected"
	PEER_CONNECTED    = "connected"
)

type PeerEvent struct {
	Peer peer.Peer
	Type string
}
