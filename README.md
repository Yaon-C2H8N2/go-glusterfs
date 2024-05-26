# go-glusterfs

A Go library for deploying and managing GlusterFS clusters. Implemented as a wrapper around the GlusterFS CLI.

## Disclaimer

This project is in early stages of development. It's a simple side-project I'm doing on my own to manage my
GlusterFS deployed on my homelab with a WebUI. It is not recommended for production use. Any feedback and contributions are more than
welcome.

## Features

- [x] Probe GlusterFS nodes
- [x] Create and delete volumes
- [x] Start volumes

## Example

```go
package main

import (
	"go-glusterfs.yaon.fr/pkg/brick"
	"go-glusterfs.yaon.fr/pkg/peer"
	"go-glusterfs.yaon.fr/pkg/volume"
)

func main() {
	//Probing nodes 2 and 3
	hostnameList := []string{"node2", "node3"}
	for _, hostname := range hostnameList {
		err := peer.PeerProbe(hostname)
		if err != nil {
			panic(err)
		}
	}

	//Creating bricks for each nodes in the pool
	peers, err := peer.ListPeers()
	var bricks []brick.Brick
	if err != nil {
		panic(err)
	}
	for _, p := range peers {
		bricks = append(bricks, brick.Brick{
			Peer: p,
			Path: "/data/brick1",
		})
	}
	
	//Creating a volume with the bricks
	v, err := volume.CreateVolume("test-volume", bricks)
	if err != nil {
        panic(err)
    }
	err = v.Start()
	if err != nil {
        panic(err)
    }
}

```