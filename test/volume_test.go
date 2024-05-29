package test

import (
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/brick"
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/peer"
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/volume"
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
	_, err = volume.CreateReplicatedVolume("testvol", bricks)
	if err != nil {
		t.Fatalf("Failed to create volume: %s", err)
	}
}

func TestListVolumes(t *testing.T) {
	v, err := volume.ListVolumes()
	if err != nil {
		t.Fatalf("Failed to list volumes: %s", err)
	}
	if len(v) == 0 {
		t.Fatalf("No volumes found. At least one volume named 'testvol' should be present.")
	}
}

func TestStartVolume(t *testing.T) {
	v, err := volume.GetVolume("testvol")
	if err != nil {
		t.Fatalf("Failed to get volume: %s", err)
	}
	err = v.Start()
	if err != nil {
		t.Fatalf("Failed to start volume: %s", err)
	}
}

func TestStopVolume(t *testing.T) {
	v, err := volume.GetVolume("testvol")
	if err != nil {
		t.Fatalf("Failed to get volume: %s", err)
	}
	err = v.Stop()
	if err != nil {
		t.Fatalf("Failed to stop volume: %s", err)
	}
}

func TestDeleteVolume(t *testing.T) {
	err := volume.DeleteVolume("testvol")
	if err != nil {
		t.Fatalf("Failed to delete volume: %s", err)
	}
}

func TestDeleteWrongVolume(t *testing.T) {
	err := volume.DeleteVolume("wrongvol")
	if err == nil {
		t.Fatalf("Deleting a non-existing volume should return an error")
	}
}
