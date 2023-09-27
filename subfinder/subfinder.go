package subfinder

import (
	"fmt"
	"bytes"
	"os/exec"

)

func ExecSubFinder(subfinderPath, domain string) (outStr string, err error) {
	args := []string{
		"-d", domain, "-silent",
	}

	fmt.Println("subfinder domain:", domain)
	cmd := exec.Command(subfinderPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()

	if err != nil {
		return "",err
	}

	outStr = string(stdout.Bytes())
	return outStr, nil
}

