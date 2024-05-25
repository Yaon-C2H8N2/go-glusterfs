package volume

import (
	"bytes"
	"errors"
	"go-glusterfs.yaon.fr/pkg/brick"
	"go-glusterfs.yaon.fr/pkg/peer"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func CreateVolume(name string, bricks []brick.Brick) (Volume, error) {
	brickString := ""
	for _, b := range bricks {
		brickString += "" + b.Peer.Hostname + ":" + b.Path + " "
	}

	cmd := exec.Command("gluster", "volume", "create", name, brickString)
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return Volume{}, err
	}
	return Volume{}, err
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
