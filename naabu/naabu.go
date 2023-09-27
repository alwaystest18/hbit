package naabu

import (
	"bytes"
	"os/exec"
	"strconv"
)

func ExecNaabu(naabuPath string, domainsFile string, ports string, rateLimit int, force bool) (outStr string, err error) {
	args := []string{
		"-list", domainsFile, "-p", ports, "-rate", strconv.Itoa(rateLimit), "-silent",
	}

	if force {
		args = append(args, "-Pn")
	}

	cmd := exec.Command(naabuPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	outStr = string(stdout.Bytes())
	return outStr, nil
}
