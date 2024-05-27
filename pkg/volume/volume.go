package volume

import (
	"bytes"
	"fmt"
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
	cmd := exec.Command("gluster", "volume", "start", v.Name, "force")
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	return err
}

func (v Volume) Stop() error {
	cmd := exec.Command("gluster", "volume", "stop", v.Name)
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	stdin, _ := cmd.StdinPipe()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return err
	}

	// Stop needs confirmation, no way to force it
	if _, err := stdin.Write([]byte("y\n")); err != nil {
		return err
	}
	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("%v: %s", err, stderr.String())
	}

	return err
}
