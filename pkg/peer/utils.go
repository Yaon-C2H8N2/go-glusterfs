package peer

import (
	"bytes"
	"os"
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
	cmd := exec.Command("gluster", "pool", "list")
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
	for _, line := range lines[1:] {
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
			hostname := filtered[1]
			// If the hostname is localhost, we replace it with the node name
			if filtered[1] == "localhost" {
				// Hacky workaround to get the proper hostname if run in a docker container
				if os.Getenv("GLUSTERFS_NODE_NAME") != "" {
					hostname = os.Getenv("GLUSTERFS_NODE_NAME")
				} else {
					hostname, _ = os.Hostname()
				}
			}

			peer := Peer{
				UUID:     filtered[0],
				Hostname: hostname,
				State:    strings.Join(filtered[2:], " "),
			}
			peers = append(peers, peer)
		}
	}

	return peers, nil
}
