package mapcidr

import (
	"fmt"
	"bytes"
	"os/exec"

)

func ExecMapCidr(mapCidrPath string, ipsFile string, mapCidrResultFile string) (outStr string, err error) {
	args := []string{
		"-cl", ipsFile, "-aggregate-approx", "-silent", "-o", mapCidrResultFile,
	}

	fmt.Println("Generate IP Range...")
	cmd := exec.Command(mapCidrPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	outStr = string(stdout.Bytes())
	return outStr, nil
}

