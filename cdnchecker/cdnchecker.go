package cdnchecker

import (
	"fmt"
	"bytes"
	"os/exec"

)

func ExecCdnChecker(cdnCheckerPath string, domainsFile string, resolversFile string, cdnCnameListFile string, noCdnDomainsFile string, noCdnIpsFile string, useCdnDomainsFile string, domainsInfoFile string) (outStr string, err error) {
	args := []string{
		"-df", domainsFile, "-cf", cdnCnameListFile, "-r", resolversFile, "-o", noCdnDomainsFile, "-oi", noCdnIpsFile, "-oc", useCdnDomainsFile, "-od", domainsInfoFile,
	}

	fmt.Println("check CDN...")
	cmd := exec.Command(cdnCheckerPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	outStr = string(stdout.Bytes())
	return outStr, nil
}

