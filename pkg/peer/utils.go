package peer

import (
	"bytes"
	"os/exec"
	"strings"
)

func PeerProbe(hostname string) error {
	cmd := exec.Command("gluster", "peer", "probe", hostname)
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	return err
}

func ListPeers() ([]Peer, error) {
	cmd := exec.Command("gluster", "peer", "status")
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return parsePeerStatus(out.String())
}

func parsePeerStatus(out string) ([]Peer, error) {
	lines := strings.Split(out, "\n")

	var peers []Peer
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var filtered []string
		for _, str := range strings.Fields(line) {
			if str != "" {
				filtered = append(filtered, str)
			}
		}
		if len(filtered) >= 3 {
			peer := Peer{
				UUID:     filtered[0],
				Hostname: filtered[1],
				State:    strings.Join(filtered[2:], " "),
			}
			peers = append(peers, peer)
		}
	}

	return peers, nil
}
