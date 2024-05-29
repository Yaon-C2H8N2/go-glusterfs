package brick

import (
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/peer"
)

type Brick struct {
	Peer peer.Peer
	Path string
}
