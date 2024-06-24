package eventListener

import (
	"fmt"
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/peer"
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/volume"
	"sort"
	"time"
)

type EventListener struct {
	isStarted      bool
	pollingTimeout int
	OnVolumeUpdate func(event VolumeEvent)
	OnPeerUpdate   func(event PeerEvent)
}

func Default() EventListener {
	return EventListener{
		isStarted:      false,
		pollingTimeout: 100,
		OnVolumeUpdate: nil,
		OnPeerUpdate:   nil,
	}
}

func (e *EventListener) Start() error {
	if e.isStarted {
		return fmt.Errorf("EventListener already started")
	}
	e.isStarted = true

	go e.listen()

	return nil
}

func (e *EventListener) Stop() error {
	if !e.isStarted {
		return fmt.Errorf("EventListener already stopped")
	}

	e.isStarted = false

	return nil
}

func (e *EventListener) SetPollingTimeout(timeout int) {
	e.pollingTimeout = timeout
}

func (e *EventListener) getIsStarted() bool {
	return e.isStarted
}

func (e *EventListener) listen() error {
	previousPoll := time.Now()
	previousVolumes, previousPeers, err := e.poll()
	if err != nil {
		return err
	}
	for e.isStarted {
		if time.Since(previousPoll).Milliseconds() > int64(e.pollingTimeout) {
			previousPoll = time.Now()
			volumes, peers, err := e.poll()
			if err != nil {
				return err
			}
			if e.OnVolumeUpdate != nil {
				go e.findVolumeUpdates(previousVolumes, volumes)
			}
			if e.OnPeerUpdate != nil {
				go e.findPeerUpdates(previousPeers, peers)
			}
			previousVolumes = volumes
			previousPeers = peers
		}
	}
	return nil
}

func (e *EventListener) poll() ([]volume.Volume, []peer.Peer, error) {
	volumes, err := volume.ListVolumes()
	if err != nil {
		return nil, nil, err
	}
	peers, err := peer.ListPeers()
	if err != nil {
		return nil, nil, err
	}
	return volumes, peers, nil
}

func (e *EventListener) findVolumeUpdates(previousVolumes []volume.Volume, volumes []volume.Volume) {
	sort.Slice(previousVolumes, func(i, j int) bool {
		return previousVolumes[i].Name < previousVolumes[j].Name
	})
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Name < volumes[j].Name
	})
	if len(previousVolumes) < len(volumes) {
		var newVol volume.Volume
		for i, v := range volumes {
			if i >= len(previousVolumes) || (i < len(previousVolumes) && v.Name != previousVolumes[i].Name) {
				newVol = v
				break
			}
		}
		e.OnVolumeUpdate(VolumeEvent{
			Type:   VOLUME_CREATE,
			Volume: newVol,
		})
	} else if len(previousVolumes) > len(volumes) {
		var oldVol volume.Volume
		for i, v := range previousVolumes {
			if i >= len(volumes) || (i < len(volumes) && v.Name != volumes[i].Name) {
				oldVol = v
				break
			}
		}
		e.OnVolumeUpdate(VolumeEvent{
			Type:   VOLUME_DELETE,
			Volume: oldVol,
		})
	} else {
		for i, v := range volumes {
			if v.Status != previousVolumes[i].Status {
				var eventType string
				if v.Status == "Started" && previousVolumes[i].Status == "Created" {
					eventType = VOLUME_START
				} else if v.Status == "Stopped" && previousVolumes[i].Status == "Started" {
					eventType = VOLUME_STOP
				}
				e.OnVolumeUpdate(VolumeEvent{
					Type:   eventType,
					Volume: v,
				})
			}
		}
	}
}

func (e *EventListener) findPeerUpdates(previousPeers []peer.Peer, peers []peer.Peer) {
	sort.Slice(previousPeers, func(i, j int) bool {
		return previousPeers[i].UUID < previousPeers[j].UUID
	})
	sort.Slice(peers, func(i, j int) bool {
		return peers[i].UUID < peers[j].UUID
	})
	if len(previousPeers) < len(peers) {
		var newPeer peer.Peer
		for i, p := range peers {
			if i >= len(previousPeers) || (i < len(previousPeers) && p.UUID != previousPeers[i].UUID) {
				newPeer = p
				break
			}
		}
		e.OnPeerUpdate(PeerEvent{
			Type: PEER_PROBED,
			Peer: newPeer,
		})
	} else if len(previousPeers) > len(peers) {
		var oldPeer peer.Peer
		for i, p := range previousPeers {
			if i >= len(peers) || (i < len(previousPeers) && p.UUID != previousPeers[i].UUID) {
				oldPeer = p
				break
			}
		}
		e.OnPeerUpdate(PeerEvent{
			Type: PEER_REMOVED,
			Peer: oldPeer,
		})
	} else {
		for i, p := range peers {
			if p.State != previousPeers[i].State {
				var eventType string
				if p.State == "Connected" && previousPeers[i].State == "Disconnected" {
					eventType = PEER_CONNECTED
				} else if p.State == "Disconnected" && previousPeers[i].State == "Connected" {
					eventType = PEER_DISCONNECTED
				}
				e.OnPeerUpdate(PeerEvent{
					Type: eventType,
					Peer: p,
				})
			}
		}
	}
}
