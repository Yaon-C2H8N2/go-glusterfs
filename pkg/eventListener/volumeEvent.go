package eventListener

import (
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/volume"
)

const (
	VOLUME_CREATE = "create"
	VOLUME_DELETE = "delete"
	VOLUME_START  = "start"
	VOLUME_STOP   = "stop"
)

type VolumeEvent struct {
	Volume volume.Volume
	Type   string
}
