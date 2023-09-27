package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 将文件内容转为字符列表
func FileContentToList(filePath string) []string {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("file open fail: " + filePath)
		return []string{""}
	}
	contentList := strings.Split(string(fileContent), "\n")
	var newList []string
	for _, element := range contentList {
		if element != "" {
			newList = append(newList, element)
		}
	}
	return newList
}

// 字符串去重去空
func UniqueStrArr(m []string) []string {
	d := make([]string, 0)
	tempMap := make(map[string]bool, len(m))
	for _, v := range m {
		if tempMap[v] == false && len(v) > 0 {
			tempMap[v] = true
			d = append(d, v)
		}
	}
	return d
}

// 判断是否为子域名
func IsSubdomain(subdomain, domain string) bool {
	if strings.HasSuffix(subdomain, "."+domain) {
		return true
	}
	return false
}

// 过滤白名单域名，跳过检测
func WhiteFilter(domainList []string, whiteList []string) []string {
	var newDomainList []string
	for _, domain := range domainList {
		contains := false
		for _, whiteStr := range whiteList {
			if strings.Contains(domain, whiteStr) {
				contains = true
				break
			}
		}
		if !contains {
			newDomainList = append(newDomainList, domain)
		}
	}
	return newDomainList
}

// 提取存在子域名的域名，比如baidu.com，当传递的域名列表中存在xxx.xxx.baidu.com时，该方法才会保留xxx.baidu.com，便于进一步使用大字典枚举
func FilterLevelDomains(domainList []string, depth int, haveSub bool) []string {
	levelDomainList := make([]string, 0)
	for _, domain := range domainList {
		domainParts := strings.Split(domain, ".")
		if haveSub {
			if len(domainParts) > depth {
				levelDomainList = append(levelDomainList, strings.Join(domainParts[len(domainParts)-depth:], "."))
			}
		} else {
			if len(domainParts) >= depth {
				levelDomainList = append(levelDomainList, strings.Join(domainParts[len(domainParts)-depth:], "."))
			}
		}
	}
	return levelDomainList
}

// 判断ip地址是否为内网ip
func IsPrivateIP(ipString string) bool {
	ip := strings.Split(ipString, ".")
	ipPart0, _ := strconv.Atoi(ip[0])
	ipPart1, _ := strconv.Atoi(ip[1])

	// Check if the IP address is in the 10.0.0.0/8 range
	if ipPart0 == 10 {
		return true
	}

	// Check if the IP address is in the 172.16.0.0/12 range
	if ipPart0 == 172 && ipPart1 >= 16 && ipPart1 <= 31 {
		return true
	}

	// Check if the IP address is in the 192.168.0.0/16 range
	if ipPart0 == 192 && ipPart1 == 168 {
		return true
	}

	// If the IP address is not in any of the private ranges, return false
	return false
}

// 把列表内容写入文件
func CreateFileWithArr(stringArr []string, fileName string) bool {
	currFile, err := os.OpenFile(fileName, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return false
	}
	for _, content := range stringArr {
		currFile.WriteString(content + "\n")
	}
	defer currFile.Close()
	return true
}

// 判断文件或目录是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func RemoveFile(path string) {
	if FileExists(path) {
		os.Remove(path)
	}
}
