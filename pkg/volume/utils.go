package volume

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/brick"
	"github.com/Yaon-C2H8N2/go-glusterfs/pkg/peer"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func CreateReplicatedVolume(name string, bricks []brick.Brick) (Volume, error) {
	var brickArgs []string
	for _, b := range bricks {
		brickArgs = append(brickArgs, b.Peer.Hostname+":"+b.Path)
	}
	args := []string{"volume", "create", name, "replica", strconv.Itoa(len(bricks))}
	args = append(args, brickArgs...)
	args = append(args, "force")

	cmd := exec.Command("gluster", args...)
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return Volume{}, err
	}
	vol, err := GetVolume(name)
	return vol, err
}

func DeleteVolume(name string) error {
	cmd := exec.Command("gluster", "volume", "delete", name)
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	stdin, _ := cmd.StdinPipe()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return err
	}

	// Deletion needs confirmation, no way to force it
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

func GetVolume(name string) (Volume, error) {
	cmd := exec.Command("gluster", "volume", "info", name)
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return Volume{}, err
	}
	vol, err := parseVolumeList(out.String())
	if err != nil {
		return Volume{}, err
	}
	if len(vol) == 0 {
		return Volume{}, errors.New("volume not found")
	}
	return vol[0], err
}

func ListVolumes() ([]Volume, error) {
	cmd := exec.Command("gluster", "volume", "info")
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return parseVolumeList(out.String())
}

func parseVolumeList(out string) ([]Volume, error) {
	lines := strings.Split(out, "\n")

	res := make(map[string]*Volume)
	peerList, _ := peer.ListPeers()
	reVolumeName := regexp.MustCompile("^Volume Name:\\s*(\\S+)")
	reType := regexp.MustCompile("^Type:\\s(\\S+)")
	reStatus := regexp.MustCompile("^Status:\\s(\\S+)")
	reNumberOfBricks := regexp.MustCompile("^Number of Bricks:\\s(\\S+)")
	reBricks := regexp.MustCompile("^Brick(\\d):\\s*(\\S+)")

	var volumeName string
	var numberOfBricks int64
	numberOfBricks = 0
	for _, line := range lines {
		isVolName, _ := regexp.MatchString(reVolumeName.String(), line)
		if isVolName {
			if numberOfBricks != 0 && int64(len(res[volumeName].Bricks)) != numberOfBricks {
				return nil, errors.New("parsing error : number of bricks mismatched number of parsed bricks")
			}
			volumeName = reVolumeName.FindStringSubmatch(line)[1]
			res[volumeName] = &Volume{Name: volumeName}
		}

		isType, _ := regexp.MatchString(reType.String(), line)
		if isType {
			res[volumeName].Type = reType.FindStringSubmatch(line)[1]
		}

		isStatus, _ := regexp.MatchString(reStatus.String(), line)
		if isStatus {
			res[volumeName].Status = reStatus.FindStringSubmatch(line)[1]
		}

		isNumberOfBricks, _ := regexp.MatchString(reNumberOfBricks.String(), line)
		if isNumberOfBricks {
			numberOfBricks, _ = strconv.ParseInt(reNumberOfBricks.FindStringSubmatch(line)[1], 10, 64)
		}

		isBrick, _ := regexp.MatchString(reBricks.String(), line)
		if isBrick {
			var brickPeer peer.Peer
			parsedBrick := strings.Split(reBricks.FindStringSubmatch(line)[2], ":")

			if len(parsedBrick) < 2 {
				return nil, errors.New("parsing error : brick couldn't be parsed")
			}

			for _, p := range peerList {
				if p.Hostname == parsedBrick[0] {
					brickPeer = p
				}
			}

			b := brick.Brick{
				Peer: brickPeer,
				Path: parsedBrick[1],
			}
			res[volumeName].Bricks = append(res[volumeName].Bricks, b)
		}
	}
	var resToArray []Volume
	for _, volume := range res {
		resToArray = append(resToArray, *volume)
	}

	return resToArray, nil
}
