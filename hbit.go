package main

import (
	"flag"
	"fmt"
	"hbit/alterx"
	"hbit/cdnchecker"
	"hbit/common"
	"hbit/hostcollision"
	"hbit/httpx"
	"hbit/mapcidr"
	"hbit/naabu"
	"hbit/shuffledns"
	"hbit/subfinder"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Global struct {
		RecursionDepth    int  `yaml:"RecursionDepth"`
		RemoveTempFile    bool `yaml:"RemoveTempFile"`
		MaxBigDictEnumNum int  `yaml:"MaxBigDictEnumNum"`
	} `yaml:"global"`
	Shuffledns struct {
		ShufflednsPath     string `yaml:"ShufflednsPath"`
		ResolversFile      string `yaml:"ResolversFile"`
		TrustResolversFile string `yaml:"TrustResolversFile"`
		BigWordListFile    string `yaml:"BigWordListFile"`
		SmallWordListFile  string `yaml:"SmallWordListFile"`
		RateLimit          int    `yaml:"RateLimit"`
		RunAsRoot          bool   `yaml:"RunAsRoot"`
	} `yaml:"shuffledns"`
	SubFinder struct {
		SubFinderPath string `yaml:"SubFinderPath"`
	} `yaml:"subfinder"`
	CdnChecker struct {
		CdnCheckerPath   string `yaml:"CdnCheckerPath"`
		ResolversFile    string `yaml:"ResolversFile"`
		CdnCnameListFile string `yaml:"CdnCnameListFile"`
	} `yaml:"cdnchecker"`
	MapCidr struct {
		MapCidrPath string `yaml:"MapCidrPath"`
	} `yaml:"mapcidr"`
	AlterX struct {
		AlterXPath string `yaml:"AlterXPath"`
		LimitNum   int    `yaml:"LimitNum"`
	} `yaml:"alterx"`
	Naabu struct {
		NaabuPath      string `yaml:"NaabuPath"`
		MiniScanPorts  string `yaml:"MiniScanPorts"`
		LargeScanPorts string `yaml:"LargeScanPorts"`
		RateLimit      int    `yaml:"RateLimit"`
	} `yaml:"naabu"`
	Httpx struct {
		HttpxPath string `yaml:"HttpxPath"`
		Threads   int    `yaml:"Threads"`
		RateLimit int    `yaml:"RateLimit"`
	} `yaml:"httpx"`
	HostCollision struct {
		HostCollisionPath string `yaml:"HostCollisionPath"`
		Threads           int    `yaml:"Threads"`
		RateLimit         int    `yaml:"RateLimit"`
		UseHostCollision  bool   `yaml:"UseHostCollision"`
	} `yaml:"hostcollision"`
}

func main() {
	df := flag.String("df", "", "[必选参数]，指定需要收集信息的域名列表文件路径，文件中每行一个域名")
	c := flag.String("c", "config/config.yaml", "非必选参数，指定配置文件路径，默认为config/config.yaml")
	i := flag.String("i", "", "非必选参数，用来导入其他方式收集到的域名，指定域名列表文件路径，文件中每行一个域名")
	wf := flag.String("wf", "", "非必选参数，用来指定不希望去枚举的子域名列表文件，比如qq的qzone.qq.com，枚举的话实际意义不大，并且会浪费很多时间")
	flag.Parse()

	//读取配置
	configFile := *c
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("读取配置文件失败 %s", err)
		os.Exit(2)
	}
	config := new(Config)
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		fmt.Printf("解析配置文件失败，请检查配置文件格式 %s", err)
		os.Exit(2)
	}

	//使用-i参数指定文件，文件如果不存在报错
	importDomainsFile := *i
	if importDomainsFile != "" {
		if !common.FileExists(importDomainsFile) {
			fmt.Printf("%s 在系统中未找到，请确认路径是否正确", importDomainsFile)
			os.Exit(2)
		}
	}

	//使用-wf参数指定文件，文件如果不存在报错，存在则转为列表
	whiteListFile := *wf
	var whiteList []string
	if whiteListFile != "" {
		if !common.FileExists(whiteListFile) {
			fmt.Printf("%s 在系统中未找到，请确认路径是否正确", whiteListFile)
			os.Exit(2)
		}
		whiteList = common.UniqueStrArr(common.FileContentToList(whiteListFile))
	}

	inputDomainFile := *df
	inputDomainsList := common.UniqueStrArr(common.FileContentToList(inputDomainFile))

	shufflednsPath := config.Shuffledns.ShufflednsPath
	subFinderPath := config.SubFinder.SubFinderPath
	cdnCheckerPath := config.CdnChecker.CdnCheckerPath
	mapCidrPath := config.MapCidr.MapCidrPath
	alterxPath := config.AlterX.AlterXPath
	naabuPath := config.Naabu.NaabuPath
	httpxPath := config.Httpx.HttpxPath
	hostCollisionPath := config.HostCollision.HostCollisionPath
	useHostCollision := config.HostCollision.UseHostCollision

	if !common.FileExists(shufflednsPath) {
		fmt.Printf("%s 在系统中未找到，请在config.yaml中配置", shufflednsPath)
		os.Exit(2)
	}

	if !common.FileExists(subFinderPath) {
		fmt.Printf("%s 在系统中未找到，请在config.yaml中配置", subFinderPath)
		os.Exit(2)
	}

	if !common.FileExists(cdnCheckerPath) {
		fmt.Printf("%s 在系统中未找到，请在config.yaml中配置", cdnCheckerPath)
		os.Exit(2)
	}

	if !common.FileExists(mapCidrPath) {
		fmt.Printf("%s 在系统中未找到，请在config.yaml中配置", mapCidrPath)
		os.Exit(2)
	}

	if !common.FileExists(alterxPath) {
		fmt.Printf("%s 在系统中未找到，请在config.yaml中配置", alterxPath)
		os.Exit(2)
	}

	if !common.FileExists(naabuPath) {
		fmt.Printf("%s 在系统中未找到，请在config.yaml中配置", naabuPath)
		os.Exit(2)
	}

	if !common.FileExists(httpxPath) {
		fmt.Printf("%s 在系统中未找到，请在config.yaml中配置", httpxPath)
		os.Exit(2)
	}

	if useHostCollision && !common.FileExists(hostCollisionPath) {
		fmt.Printf("%s 在系统中未找到，请在config.yaml中配置", hostCollisionPath)
		os.Exit(2)
	}

	//创建结果保存根目录，目录名为当前日期
	reportRootDir := "reports/"
	if !common.FileExists(reportRootDir) {
		err := os.MkdirAll(reportRootDir, os.ModePerm)
		if err != nil {
			fmt.Println("结果保存目录创建失败，请检查权限")
			os.Exit(2)
		}
	}

	//创建结果保存根目录，目录名为当前日期
	reportDateDir := reportRootDir + time.Now().Format("20060102") + "/"
	if !common.FileExists(reportDateDir) {
		err := os.MkdirAll(reportDateDir, os.ModePerm)
		if err != nil {
			fmt.Println("结果保存二级目录创建失败，请检查权限")
			os.Exit(2)
		}
	}

	//按域名创建结果保存目录
	for _, domain := range inputDomainsList {
		reportDir := reportDateDir + domain + "/"
		if !common.FileExists(reportDir) {
			err := os.MkdirAll(reportDir, os.ModePerm)
			if err != nil {
				fmt.Println("域名对应结果目录创建失败，请检查权限")
				os.Exit(2)
			}
		}

		//定义临时域名列表
		var tempDomainList []string

		//将其他途径获取到的子域名添加至临时域名列表
		if importDomainsFile != "" {
			importDomainsList := common.UniqueStrArr(common.FileContentToList(importDomainsFile))
			tempDomainList = append(tempDomainList, importDomainsList...)
		}

		//使用shuffledns枚举相对二级域名
		fmt.Println("enumerate domain: " + domain)
		shufflednsResult, _ := shuffledns.ExecShuffleDns(
			shufflednsPath,
			domain,
			config.Shuffledns.ResolversFile,
			config.Shuffledns.BigWordListFile,
			config.Shuffledns.RateLimit,
			config.Shuffledns.RunAsRoot)

		//结果添加至临时域名列表
		shufflednsDomainList := strings.Split(shufflednsResult, "\n")
		tempDomainList = append(tempDomainList, shufflednsDomainList...)

		//使用subfinder查找子域名，结果添加至临时域名列表
		subfinderResult, _ := subfinder.ExecSubFinder(config.SubFinder.SubFinderPath, domain)
		subfinderResultList := strings.Split(subfinderResult, "\n")
		tempDomainList = append(tempDomainList, subfinderResultList...)

		//定义临时子域名列表，用来排除导入或收集结果中不属于该域名的子域名
		var tempSubDomainList []string

		//判断是否属于检测域名的子域（subfinder或者导入域名有可能为其他域下的子域名）
		for _, tempSubDomain := range tempDomainList {
			if common.IsSubdomain(tempSubDomain, domain) {
				tempSubDomainList = append(tempSubDomainList, tempSubDomain)
			}
		}

		tempDomainList = tempSubDomainList

		//过滤掉白名单域名
		if len(whiteList) > 0 {
			tempDomainList = common.WhiteFilter(tempDomainList, whiteList)
		}

		//临时域名列表去重，里面包含失效域名
		tempDomainList = common.UniqueStrArr(tempDomainList)

		//临时域名列表写入文件，便于shuffledns解析域名
		tempDomainsFile := reportDir + "temp_Domains_" + domain
		if len(tempDomainList) > 0 {
			if !common.CreateFileWithArr(tempDomainList, tempDomainsFile) {
				fmt.Printf("%s File create fail", tempDomainsFile)
				os.Exit(2)
			}
		} else {
			fmt.Println("Subdomain not found\n")
			continue
		}

		//使用shuffledns解析域名，利于大dns列表，速度快，但准确率无法保证，后续通过可信dns二次验证
		tempVaildDomainsFile := reportDir + "temp_Vaild_Domains_" + domain
		shuffledns.ExecShuffleDnsResolv(
			shufflednsPath,
			domain,
			tempDomainsFile,
			config.Shuffledns.ResolversFile,
			config.Shuffledns.RateLimit,
			tempVaildDomainsFile,
			config.Shuffledns.RunAsRoot)

		tempLevelAllVaildDomainsFile := reportDir + "temp_LevelAll_Vaild_Domains_" + domain

		//对shuffledns验证的域名进行可信dns的二次验证，确保域名真实存在
		shuffledns.ExecShuffleDnsResolv(
			shufflednsPath,
			domain,
			tempVaildDomainsFile,
			config.Shuffledns.TrustResolversFile,
			config.Shuffledns.RateLimit,
			tempLevelAllVaildDomainsFile,
			config.Shuffledns.RunAsRoot)

		var domainList []string
		domainList = common.UniqueStrArr(common.FileContentToList(tempLevelAllVaildDomainsFile))

		//域名递归枚举
		lenDomain := len(strings.Split(domain, "."))
		depth := config.Global.RecursionDepth
		maxBigDictEnumNum := config.Global.MaxBigDictEnumNum
		if depth > 1 {
			enumDomainDepth := lenDomain + depth - 1
			for i := lenDomain + 1; i <= enumDomainDepth; i++ {
				levelDomains := common.UniqueStrArr(common.FilterLevelDomains(domainList, i, false))
				tempLevelDomainsFile := reportDir + "temp_level" + strconv.Itoa(i) + "_Domains_" + domain
				if !common.CreateFileWithArr(levelDomains, tempLevelDomainsFile) {
					fmt.Println("temp level domains file create fail")
				}
				//如果存在子域名数量大于阈值，则先对每个子域名小字典枚举，然后对存在下级子域名的域名做大字典枚举，避免对所有子域名大字典枚举增加耗时
				if len(levelDomains) > maxBigDictEnumNum {
					var miniCount int = len(levelDomains)
					var tempMiniLevelDomainsList []string
					//小字典枚举
					for _, levelDomain := range levelDomains {
						miniCount = miniCount - 1
						fmt.Printf("mini dict enumerate subdomain: %s   number of remaining domains: %d \n", levelDomain, miniCount)
						tempLevelDomainResult, _ := shuffledns.ExecShuffleDns(
							shufflednsPath,
							levelDomain,
							config.Shuffledns.ResolversFile,
							config.Shuffledns.SmallWordListFile,
							config.Shuffledns.RateLimit,
							config.Shuffledns.RunAsRoot)

						tempMiniLevelDomainsList = append(tempMiniLevelDomainsList, strings.Split(tempLevelDomainResult, "\n")...)

					}
					//小字典枚举域名保存到临时文件，便于后续二次验证使用
					tempMiniAllDomainFile := reportDir + "temp_level" + strconv.Itoa(i) + "_Mini_All_Domains_" + domain
					common.CreateFileWithArr(tempMiniLevelDomainsList, tempMiniAllDomainFile)

					//对小字典枚举域名使用可信dns二次验证
					tempMiniVaildDomainFile := reportDir + "temp_level" + strconv.Itoa(i) + "_Mini_Vaild_Domains_" + domain
					shuffledns.ExecShuffleDnsResolv(
						shufflednsPath,
						domain,
						tempMiniAllDomainFile,
						config.Shuffledns.TrustResolversFile,
						config.Shuffledns.RateLimit,
						tempMiniVaildDomainFile,
						config.Shuffledns.RunAsRoot)

					domainList = append(domainList, common.FileContentToList(tempMiniVaildDomainFile)...)

					//删除临时文件
					common.RemoveFile(tempMiniAllDomainFile)
					common.RemoveFile(tempMiniVaildDomainFile)

					levelDomains = common.UniqueStrArr(common.FilterLevelDomains(domainList, i, true))
					fmt.Println("big dict enum: \n")
					//对存在下级子域名的域名做大字典枚举
					var bigCount int = len(levelDomains)
					var tempBigLevelDomainsList []string //大字典枚举域名结果临时保存列表
					for _, levelDomain := range levelDomains {
						bigCount = bigCount - 1
						fmt.Printf("big dict enumerate subdomain: %s   number of remaining domains: %d \n", levelDomain, bigCount)
						tempLevelDomainResult, _ := shuffledns.ExecShuffleDns(
							shufflednsPath,
							levelDomain,
							config.Shuffledns.ResolversFile,
							config.Shuffledns.BigWordListFile,
							config.Shuffledns.RateLimit,
							config.Shuffledns.RunAsRoot)

						tempBigLevelDomainsList = append(tempBigLevelDomainsList, strings.Split(tempLevelDomainResult, "\n")...)

					}
					//大字典枚举域名保存到临时文件，便于后续二次验证使用
					tempBigAllDomainFile := reportDir + "temp_level" + strconv.Itoa(i) + "_Big_All_Domains_" + domain
					common.CreateFileWithArr(tempBigLevelDomainsList, tempBigAllDomainFile)

					//对大字典枚举域名使用可信dns二次验证
					tempBigVaildDomainFile := reportDir + "temp_level" + strconv.Itoa(i) + "_Big_Vaild_Domains_" + domain
					shuffledns.ExecShuffleDnsResolv(
						shufflednsPath,
						domain,
						tempBigAllDomainFile,
						config.Shuffledns.TrustResolversFile,
						config.Shuffledns.RateLimit,
						tempBigVaildDomainFile,
						config.Shuffledns.RunAsRoot)

					domainList = append(domainList, common.FileContentToList(tempBigVaildDomainFile)...)

					//删除临时文件
					common.RemoveFile(tempBigAllDomainFile)
					common.RemoveFile(tempBigVaildDomainFile)

				} else if len(levelDomains) > 0 {
					//全部子域名大字典枚举
					var bigCount int = len(levelDomains)
					var tempBigLevelDomainsList []string //大字典枚举域名结果临时保存列表
					for _, levelDomain := range levelDomains {
						bigCount = bigCount - 1
						fmt.Printf("big dict enumerate subdomain: %s   number of remaining domains: %d \n", levelDomain, bigCount)
						tempLevelDomainResult, _ := shuffledns.ExecShuffleDns(
							shufflednsPath,
							levelDomain,
							config.Shuffledns.ResolversFile,
							config.Shuffledns.BigWordListFile,
							config.Shuffledns.RateLimit,
							config.Shuffledns.RunAsRoot)

						tempBigLevelDomainsList = append(tempBigLevelDomainsList, strings.Split(tempLevelDomainResult, "\n")...)
					}

					//大字典枚举域名保存到临时文件，便于后续二次验证使用
					tempBigAllDomainFile := reportDir + "temp_level" + strconv.Itoa(i) + "_Big_All_Domains_" + domain
					common.CreateFileWithArr(tempBigLevelDomainsList, tempBigAllDomainFile)

					//对大字典枚举域名使用可信dns二次验证
					tempBigVaildDomainFile := reportDir + "temp_level" + strconv.Itoa(i) + "_Big_Vaild_Domains_" + domain
					shuffledns.ExecShuffleDnsResolv(
						shufflednsPath,
						domain,
						tempBigAllDomainFile,
						config.Shuffledns.TrustResolversFile,
						config.Shuffledns.RateLimit,
						tempBigVaildDomainFile,
						config.Shuffledns.RunAsRoot)

					domainList = append(domainList, common.FileContentToList(tempBigVaildDomainFile)...)

					//删除临时文件
					common.RemoveFile(tempBigAllDomainFile)
					common.RemoveFile(tempBigVaildDomainFile)
				}
				common.RemoveFile(tempLevelDomainsFile)
			}
		}

		tempInputAlterxFile := reportDir + "temp_input_alterx_" + domain   //供alterx生成字典的域名列表
		tempOutputAlterxFile := reportDir + "temp_output_alterx_" + domain //alterx输出文件

		//如果供Alterx生成字典的域名列表保存失败，则退出程序
		if !common.CreateFileWithArr(domainList, tempInputAlterxFile) {
			fmt.Println("Alterx input file create fail")
			os.Exit(2)
		}

		//使用alterx生成域名字典
		alterx.ExecAlterX(alterxPath, tempInputAlterxFile, config.AlterX.LimitNum, tempOutputAlterxFile)

		//对alterx生成的域名字典解析，找出存在域名
		tempResolvAlterxFile := reportDir + "temp_Resolv_Alterx_Domains_" + domain
		if common.FileExists(tempOutputAlterxFile) {
			//use shuffledns to resolv domains generated by alterx
			shuffledns.ExecShuffleDnsResolv(
				shufflednsPath,
				domain,
				tempOutputAlterxFile,
				config.Shuffledns.ResolversFile,
				config.Shuffledns.RateLimit,
				tempResolvAlterxFile,
				config.Shuffledns.RunAsRoot)

			domainList = common.UniqueStrArr(append(domainList, common.FileContentToList(tempResolvAlterxFile)...))
		} else {
			fmt.Printf("Alterx output file not found")
		}

		//域名汇总列表写入文件
		domainsFile := reportDir + "all_Domains_" + domain
		if !common.CreateFileWithArr(domainList, domainsFile) {
			fmt.Println("allDomains file create fail")
		}

		//调用cdnchecker筛选未使用cdn的ip
		noCdnDomainsFile := reportDir + "temp_" + domain + "_nocdn_domains"
		useCdnDomainsFile := reportDir + "temp_" + domain + "_usecdn_domains"
		noCdnIpsFile := reportDir + "temp_" + domain + "_nocdn_ips"
		domainsInfoFile := reportDir + "domains_info_" + domain
		cdncheckerResult, err := cdnchecker.ExecCdnChecker(
			config.CdnChecker.CdnCheckerPath,
			domainsFile,
			config.CdnChecker.ResolversFile,
			config.CdnChecker.CdnCnameListFile,
			noCdnDomainsFile,
			noCdnIpsFile,
			useCdnDomainsFile,
			domainsInfoFile)
		if err != nil {
			fmt.Println("cdnChecker exec fail!")
			os.Exit(2)
		}
		fmt.Println(cdncheckerResult)

		//调用mapcidr生成ip范围
		mapCidrResultFile := reportDir + "ip_range_" + domain
		mapCidrResult, _ := mapcidr.ExecMapCidr(mapCidrPath, noCdnIpsFile, mapCidrResultFile)
		fmt.Println(mapCidrResult)

		//调用naabu扫描全域名80,443端口
		var domainPortsList []string
		miniScanPorts := config.Naabu.MiniScanPorts
		fmt.Println("check mini range ports...")
		naabuMiniScanResult, _ := naabu.ExecNaabu(naabuPath, domainsFile, miniScanPorts, config.Naabu.RateLimit, true)
		domainPortsList = common.UniqueStrArr(strings.Split(naabuMiniScanResult, "\n"))

		//调用naabu扫描ip段自定义范围端口
		largeScanPorts := config.Naabu.LargeScanPorts
		fmt.Println("check large range ports...")
		naabuIpLargeScanResult, _ := naabu.ExecNaabu(naabuPath, mapCidrResultFile, largeScanPorts, config.Naabu.RateLimit, false)
		naabuIpLargePortsList := strings.Split(naabuIpLargeScanResult, "\n")
		domainPortsList = common.UniqueStrArr(append(domainPortsList, naabuIpLargePortsList...))

		//内网域名另存，做host碰撞
		var privateDomainsList []string
		domainsInfoList := common.UniqueStrArr(common.FileContentToList(domainsInfoFile))
		if len(domainsInfoList) > 0 {
			for _, domainInfo := range domainsInfoList {
				domainInfoList := strings.Split(domainInfo, ":")
				hostName := domainInfoList[0]
				ipString := domainInfoList[1]
				if common.IsPrivateIP(ipString) {
					privateDomainsList = append(privateDomainsList, hostName)
				}
			}
			privateDomainsList = common.UniqueStrArr(privateDomainsList)
		}
		privateDomainsFile := reportDir + "private_domains_" + domain
		if len(privateDomainsList) > 0 {
			common.CreateFileWithArr(privateDomainsList, privateDomainsFile)
		}

		//把开放端口列表写入文件
		domainOpenPortsFile := reportDir + "open_ports_" + domain
		common.CreateFileWithArr(domainPortsList, domainOpenPortsFile)

		//调用httpx验证站点存活情况
		allSiteFile := reportDir + "all_Sites_" + domain
		if common.FileExists(domainOpenPortsFile) {
			httpx.ExecHttpx(
				httpxPath,
				domainOpenPortsFile,
				config.Httpx.Threads,
				config.Httpx.RateLimit,
				allSiteFile)
		}

		//调用hostCollision对解析到内网ip的域名进行host碰撞
		if common.FileExists(allSiteFile) && common.FileExists(privateDomainsFile) && useHostCollision {
			resultHostCollisionFile := reportDir + "Host_Collision_" + domain
			hostcollision.ExecHostCollision(
				hostCollisionPath,
				privateDomainsFile,
				allSiteFile,
				config.HostCollision.Threads,
				config.HostCollision.RateLimit,
				resultHostCollisionFile)
		}

		//删除临时文件
		if config.Global.RemoveTempFile {
			common.RemoveFile(tempDomainsFile)
			common.RemoveFile(tempInputAlterxFile)
			common.RemoveFile(tempOutputAlterxFile)
			common.RemoveFile(noCdnDomainsFile)
			common.RemoveFile(useCdnDomainsFile)
			common.RemoveFile(noCdnIpsFile)
			common.RemoveFile(tempResolvAlterxFile)
			common.RemoveFile(tempVaildDomainsFile)
			common.RemoveFile(tempLevelAllVaildDomainsFile)
		}
	}
}
