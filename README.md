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

**首先需要在hbit目录创建domain.txt文件，将待收集信息的域名写入，每行一个分隔**

```
docker run -it --rm -v ./reports:/tools/hbit/reports -v ./domain.txt:/tools/hbit/domain.txt -v ./provider-config.yaml:/root/.config/subfinder/provider-config.yaml -v ./config/config.yaml:/tools/hbit/config/config.yaml hbit /bin/bash -c "cd /tools/hbit && ./hbit -df domain.txt"
```



### 手工部署方式启动

```
./hbit -df domain.txt   //domain.txt根据实际情况替换
```



### 重要配置

部分配置对结果会产生较大影响，默认配置取一个覆盖率与时间的平衡，针对特定目标，大家可以根据配置自行优化来达到最理想的效果

- provider-config.yaml  此文件为subfinder的api key配置文件，原理就是通过像大家熟知的fofa、zoomeye这类站点获取目标资产，未配置api key的情况下获取数据量非常有限，因此大家一定要尽量多的配置，配置格式可参考https://github.com/projectdiscovery/subfinder#post-installation-instructions 

- 以下介绍配置项均在config/config.yaml中，**仅为对结果影响较大的点，注意非全部配置**

  ```
  global:
    RecursionDepth: 3        #递归深度，默认为3，即域名为a.com，递归到x.x.x.a.com，也就是我们常说的四级域名，根据目标实际情况大家可以把此值继续增大，但大量的递归也意味着速度更慢
    MaxBigDictEnumNum: 20   #如果存在子域名数量大于设定值，则先对每个子域名小字典枚举，然后对存在下级子域名的域名做大字典枚举，避免对所有子域名使用大字典枚举增加耗时。小字典毕竟覆盖有限，大家如果为了更全的结果接受多花费一些时间的话，可以把此项值调大
  
  shuffledns:
    ResolversFile: '/tools/dnsVerifier/resolvers.txt'    #dns服务器文件，更多优质的dns服务器对域名枚举的速度起到非常大的作用，docker中会使用我的另一款自研程序https://github.com/alwaystest18/dnsVerifier筛选，大家有更优质的dns资源可以自行替换
    TrustResolversFile: '/tools/hbit/trustDns'         #可信dns服务器文件，部分dns服务器会产生一些垃圾数据，因此跑出来的域名要经过可信dns过滤一遍，大家有其他完全可信dns服务器可以添加到此文件，节省过滤的时间
    BigWordListFile: '/tools/dict/domain_20000'            #域名枚举大字典，好的字典对结果的影响不言而喻，比较小的目标我都用40w行的字典跑
    SmallWordListFile: '/tools/dict/small_domain_dict'    #域名枚举小字典，道理同上，大家根据实际情况替换字典文件
    RateLimit: 10000               #家庭宽带建议调整至200，否则可能因流量过大引起断网，网络好的可以继续往大调整
  
  
  cdnchecker: 
    ResolversFile: '/tools/cdnChecker/resolvers_cn.txt'     #dns服务器文件,最好是国内节点，否则会影响cdn识别准确率，毕竟从国外解析国内cdn的域名，ip段可能是相同的，就会产生误报
    CdnCnameListFile: '/tools/cdnChecker/cdn_cname'          #cdncname文件，已知cdn的cname大家都可以自行往里加
  
  
  alterx:
    LimitNum: 1000000                                  #alterx生成子域名字典最大行数，对于一些比较大的目标（子域名1000+），可以适当调大此值，不过行数越多也代表时间越久
  
  
  naabu:
    RateLimit: 1000                              #速率限制，家庭宽带建议调整至200，否则可能因流量过大引起断网，网络好的可以继续往大调整
    MiniScanPorts: '80,443'                   #端口扫描最小范围，全部站点使用此范围扫描，对于用了cdn的，往往就是这两个端口，不会有22 3389这种，所以配置太多只会影响速度
    LargeScanPorts: '80,81...60443'   #端口扫描大范围，仅用于未使用cdn站点，如果设置范围过大，比如1-65535，虽然会增加资产检出率，但严重影响速度
  
  
  ```

  
