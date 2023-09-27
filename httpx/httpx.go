package httpx

import (
	"fmt"
	"bytes"
	"os/exec"
	"strconv"

)

func ExecHttpx(httpxPath string, domainsFile string, threads int, rateLimit int, outPutFile string) (outStr string, err error) {
	args := []string{
		"-list", domainsFile, "-t", strconv.Itoa(threads), "-rl", strconv.Itoa(rateLimit), "-silent", "-o", outPutFile,
	}

	fmt.Println("use httpx to verify site alive...")
	cmd := exec.Command(httpxPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	outStr = string(stdout.Bytes())
	return outStr, nil
}

