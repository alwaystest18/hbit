package shuffledns

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
)

func ExecShuffleDns(shufflednsPath string, domain string, resolversFile string, wordListFile string, rateLimit int, runAsRoot bool) (outStr string, err error) {
	args := []string{
		"-d", domain, "-w", wordListFile, "-r", resolversFile, "-t", strconv.Itoa(rateLimit),
		"-silent",
	}
	if runAsRoot {
		args = append(args, "-mcmd", "--root")
	}

	cmd := exec.Command(shufflednsPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	outStr = string(stdout.Bytes())
	return outStr, nil
}

func ExecShuffleDnsResolv(shufflednsPath string, domain string, domainListFile string, resolversFile string, rateLimit int, outputFile string, runAsRoot bool) (outStr string, err error) {
	args := []string{
		"-d", domain, "-list", domainListFile, "-r", resolversFile, "-t", strconv.Itoa(rateLimit), "-o", outputFile,
		"-silent",
	}
	if runAsRoot {
		args = append(args, "-mcmd", "--root")
	}

	fmt.Println("resolv subdomain:", domain)
	cmd := exec.Command(shufflednsPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	outStr = string(stdout.Bytes())
	return outStr, nil
}
