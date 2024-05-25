package brick

import (
	"go-glusterfs.yaon.fr/pkg/peer"
)

type Brick struct {
	Peer peer.Peer
	Path string
}
