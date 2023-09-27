package alterx

import (
	"fmt"
	"bytes"
	"os/exec"
	"strconv"
)

func ExecAlterX(alterxPath string, domainsFile string, limitNum int, alterXResultFile string) (outStr string, err error) {
	args := []string{
		"-list", domainsFile, "-en", "-limit", strconv.Itoa(limitNum), "-silent", "-o", alterXResultFile,
	}

	fmt.Println("Use AlterX Generate Domain List...")
	cmd := exec.Command(alterxPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	outStr = string(stdout.Bytes())
	return outStr, nil
}

