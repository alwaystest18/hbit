# hbit

一款用于安全测试中信息收集的自动化工具

## 介绍

一款信息收集的缝合工具，用到了以下工具

[shuffledns](https://github.com/projectdiscovery/shuffledns)  子域名枚举，根据配置可递归枚举

[subfinder](https://github.com/projectdiscovery/subfinder)    被动收集子域名

[alterx](https://github.com/projectdiscovery/alterx)  根据已知子域名规律生成字典

[cdnChecker](https://github.com/alwaystest18/cdnChecker)  cdn识别

[mapcidr](https://github.com/projectdiscovery/mapcidr)   根据ip地址自动生成网段格式

[naabu](https://github.com/projectdiscovery/naabu)   端口扫描

[httpx](https://github.com/projectdiscovery/httpx)   http服务存活检测

[hostCollision](https://github.com/alwaystest18/hostCollision)  host碰撞



程序输出以下文件：

all_Domains_{domain}  对应域名下的全部子域名   示例：sub.test.com

all_Sites_{domain}   对应域名下的全部http站点   示例：https://sub.test.com:8443

domains_info_{domain}   子域名及对应ip地址，可根据ip地址反查相关域名，便于确认资产  示例：sub.test.com:1.1.1.1

ip_range_{domain}  根据已知资产生成的ip段范围  示例：1.1.1.0/24

open_ports_{domain}  对应域名下的资产开放端口信息  示例：sub.test.com:3389

private_domains_{domain}  如果存在解析ip为内网域名会产生此文件，便于进一步做host碰撞  示例：private.test.com

Host_Collision_{domain}  如果存在host碰撞成功的内网域名会产生此文件  示例：url:http://1.1.1.1  host:private.test.com  title:[test]  Length: 666



### 程序特点

- cdn识别：用到了自研程序cdnChecker，通过cname+多节点解析域名结合的方式判断是否使用cdn，相比单独依靠cname去判断cdn，准确率有所提升;

- Host碰撞：用到了自研程序hostCollision，对解析ip为内网的域名做自动化的host碰撞，尽早发现外网可访问的内网资产;

- 输出简洁：输出结果无需额外处理即可直接与大部分工具联动;

- 可兼容其他程序结果深入利用：有很多优秀的信息收集工具，如[OneForAll](https://github.com/shmilylty/OneForAll) 收集方式不完全相同，因此会存在一些本工具检测不到而其他工具可检测到的资产，通过-i参数可导入已知子域名（其他工具或自己人工收集到的），本工具会结合已知域名规律生成字典进一步去枚举子域名及后续的扫端口，检测cdn等操作;

- 支持同时对多个域名的信息收集：大型攻防演练往往会有多个目标，此功能可帮助使用者提升效率。



### 程序缺点

相比其他信息收集类的工具，本工具最大的缺点就是慢，经过测试大部分域名的信息收集流程需要（10min~40min），部分超大企业，可能需要几个小时，影响运行速度有因素主要有以下几点：

- 枚举子域名阶段：存在子域名较多；域名字典较大；可用的dns服务器较少；
- 端口扫描阶段：ip较多；端口范围设置较大
- 验证存活http服务阶段：开放端口较多（部分主机扫描端口全部开放，使用httpx验证存活严重影响速度）
- host碰撞阶段：内网域名较多；存活http服务较多；



## 安装

<details>
<summary><b> docker（强烈推荐）</b></summary>

程序依赖工具较多，且涉及大量配置，因此强烈推荐使用docker一键部署

```
git clone https://github.com/alwaystest18/hbit.git
cd hbit
docker build -t hbit .
```
</details>

<details>
<summary><b> 手工</b></summary>

这里以centos7举例，详细程序安装方式可参考对应程序的github主页，注意根据部署的实际情况修改config/config.yaml文件

**部署massdns**

```
mkdir /tools
wget https://github.com/blechschmidt/massdns/archive/refs/tags/v1.0.0.tar.gz
tar zvxf v1.0.0.tar.gz
cd massdns-1.0.0/
make
ln -s /tools/massdns-1.0.0/bin/massdns /usr/bin/massdns
```

**部署subfinder**

```
go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
```

**部署shuffledns**

```
go install -v github.com/projectdiscovery/shuffledns/cmd/shuffledns@latest
```

**部署mapcidr**

```
go install -v github.com/projectdiscovery/mapcidr/cmd/mapcidr@latest
```

**部署alterx**

```
go install github.com/projectdiscovery/alterx/cmd/alterx@latest
```

**部署naabu**

```
go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@latest
```

**部署httpx**

```
go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest
```

**部署cdnChecker**

```
git clone https://github.com/alwaystest18/cdnChecker.git
cd cdnChecker/
go install
go build cdnChecker.go
```

**部署hostCollision**

```
git clone https://github.com/alwaystest18/hostCollision.git
cd hostCollision/
go install
go build hostCollision.go
```

**部署hbit**
```
git clone https://github.com/alwaystest18/hbit.git
cd hbit/
go install
go build hbit.go
```
</details>




## 使用

### 参数说明

```
Usage of ./hbit:
  -c string  //非必选参数，指定配置文件路径，默认为config/config.yaml
        config file path (default "config/config.yaml")
  -df string    //必选参数，指定需要收集信息的域名列表文件路径，文件中每行一个域名
        domain list file
  -i string   //非必选参数，用来导入其他方式收集到的域名，指定域名列表文件路径，文件中每行一个域名
        import domains list file
  -wf string   //非必选参数，跳过指定域名资产，用来指定不希望去枚举的子域名列表文件，比如qq的qzone.qq.com，枚举的话实际意义不大，并且会浪费很多时间
        white list file
```

### docker方式启动

subfinder依赖的api key可参考https://github.com/projectdiscovery/subfinder#post-installation-instructions 填写至`provider-config.yaml`文件中，如没有文件留空即可

```
docker run -it --rm -v ./reports:/tools/hbit/reports -v ./domain.txt:/tools/hbit/domain.txt -v ./provider-config.yaml:/root/.config/subfinder/provider-config.yaml -v ./config/config.yaml:/tools/hbit/config/config.yaml hbit /bin/bash -c "cd /tools/hbit && ./hbit -df domain.txt"
```



### 手工部署方式启动

```
./hbit -df domain.txt   //domain.txt根据实际情况替换
```

