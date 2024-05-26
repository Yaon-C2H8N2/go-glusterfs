package test

import (
	"go-glusterfs.yaon.fr/pkg/brick"
	"go-glusterfs.yaon.fr/pkg/peer"
	"go-glusterfs.yaon.fr/pkg/volume"
	"testing"
)

func TestCreateVolume(t *testing.T) {
	peers, err := peer.ListPeers()
	if err != nil {
		t.Fatalf("Failed to list peers: %s", err)
	}
	var bricks []brick.Brick
	for _, p := range peers {
		b := brick.Brick{Peer: p, Path: "/mnt/brick1/brick"}
		bricks = append(bricks, b)
	}
	vol, err := volume.CreateVolume("testvol", bricks)
	if err != nil {
		t.Fatalf("Failed to create volume: %s", err)
	}
	t.Logf("Created volume: %+v", vol)
}

func TestListVolumes(t *testing.T) {
	v, err := volume.ListVolumes()
	if err != nil {
		t.Fatalf("Failed to list volumes: %s", err)
	}
	t.Logf("Volumes: %+v", v)
}
