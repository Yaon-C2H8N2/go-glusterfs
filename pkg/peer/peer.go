package peer

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Peer struct {
	UUID     string
	Hostname string
	State    string
}

func (p *Peer) Detach() error {
	cmd := exec.Command("gluster", "peer", "detach", p.Hostname, "force")
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
