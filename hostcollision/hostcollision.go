package hostcollision

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
)

func ExecHostCollision(hostCollisionPath string, hostsFile string, sitesFile string, threads int, rateLimit int, hostCollisionFile string) (outStr string, err error) {
	args := []string{
		"-df", hostsFile, "-uf", sitesFile, "-t", strconv.Itoa(threads), "-r", strconv.Itoa(rateLimit), "-o", hostCollisionFile,
	}

	fmt.Println("host collision...")
	cmd := exec.Command(hostCollisionPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	outStr = string(stdout.Bytes())
	return outStr, nil
}
