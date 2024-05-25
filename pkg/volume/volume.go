package volume

import (
	"bytes"
	"go-glusterfs.yaon.fr/pkg/brick"
	"os/exec"
)

type Volume struct {
	Name   string
	Type   string
	Status string
	Bricks []brick.Brick
}

func (v Volume) Start() error {
	cmd := exec.Command("gluster", "volume", "start", v.Name)
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	return err
}
