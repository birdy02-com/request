package request

import (
	"fmt"
	"net"
	"strings"
)

// CheckCDN 检测是否为CDN
func CheckCDN(ipv4, cname string, header map[string]string) IsCDN {
	var result IsCDN
	if ipv4 != "" {
		return CheckCdnIPv4(ipv4)
	}
	if cname != "" {
		return CheckCNAME(cname)
	}
	if header != nil && len(header) > 0 {
		return CheckHeaderCDN(header)
	}
	return result
}

// IsCDN 判断是否CDN结构
type IsCDN struct {
	Is   bool   `json:"isCdn"`
	Name string `json:"name"`
}

// CdnIpList CDN IP列表
var CdnIpList = []struct {
	Name    string   `yaml:"name"`
	IpCider []string `yaml:"ip_cidr"`
}{
	{
		Name: "百度云加速",
		IpCider: []string{
			"101.227.206.0/24",
			"101.227.207.0/24",
			"101.69.175.0/24",
			"111.132.134.0/24",
			"111.174.61.0/24",
			"111.174.63.0/24",
			"111.32.134.0/24",
			"111.32.135.0/24",
			"111.32.136.0/24",
			"111.63.67.0/24",
			"111.63.68.0/24",
			"112.25.89.0/24",
			"112.25.90.0/24",
			"112.25.91.0/24",
			"112.29.157.0/24",
			"112.29.158.0/24",
			"112.29.159.0/24",
			"113.207.100.0/24",
			"113.207.101.0/24",
			"113.207.102.0/24",
			"115.231.186.0/24",
			"115.231.187.0/24",
			"116.31.126.0/24",
			"116.31.127.0/24",
			"117.147.214.0/24",
			"117.147.215.0/24",
			"117.27.149.0/24",
			"117.34.13.0/24",
			"117.34.14.0/24",
			"117.34.28.0/24",
			"117.34.60.0/24",
			"117.34.61.0/24",
			"117.34.62.0/24",
			"119.147.134.0/24",
			"119.167.246.0/24",
			"119.188.132.0/24",
			"119.188.14.0/24",
			"119.188.9.0/24",
			"119.188.97.0/24",
			"119.84.1.0/24",
			"119.84.92.0/24",
			"119.84.93.0/24",
			"122.190.1.0/24",
			"122.190.2.0/24",
			"122.190.3.0/24",
			"122.246.5.0/24",
			"124.95.168.128/25",
			"124.95.188.0/24",
			"124.95.191.0/24",
			"125.39.174.0/24",
			"125.39.238.0/24",
			"125.39.239.0/24",
			"14.17.71.0/24",
			"150.138.149.0/24",
			"150.138.150.0/24",
			"150.138.151.0/24",
			"157.255.24.0/24",
			"157.255.25.0/24",
			"157.255.26.0/24",
			"180.163.113.0/24",
			"180.163.153.0/24",
			"180.163.154.0/24",
			"180.163.188.0/24",
			"180.163.189.0/24",
			"183.232.51.0/24",
			"183.232.53.0/24",
			"183.60.235.0/24",
			"183.61.177.0/24",
			"183.61.190.0/24",
			"183.61.236.0/24",
			"219.159.84.0/24",
			"220.170.184.0/24",
			"220.170.185.0/24",
			"220.170.186.0/24",
			"220.195.21.0/25",
			"220.195.22.0/24",
			"221.178.56.0/24",
			"221.178.57.0/24",
			"221.178.58.0/26",
			"222.216.190.0/24",
			"42.236.7.128/26",
			"42.236.7.64/27",
			"42.236.93.0/24",
			"42.236.94.0/24",
			"42.81.6.0/24",
			"42.81.8.0/24",
			"58.211.137.0/24",
			"58.211.2.0/24",
			"59.51.81.128/25",
			"60.217.232.0/24",
			"61.155.149.0/24",
			"61.155.165.0/24",
			"61.156.149.0/24",
			"61.182.136.0/24",
			"61.182.137.0/24",
			"61.241.118.0/24",
		},
	},
	{
		Name: "加速乐CDN",
		IpCider: []string{
			"113.107.238.0/24",
			"106.42.25.0/24",
			"183.222.96.0/24",
			"117.21.219.0/24",
			"116.55.250.0/24",
			"111.202.98.0/24",
			"111.13.147.0/24",
			"122.228.238.0/24",
			"58.58.81.0/24",
			"1.31.128.0/24",
			"123.155.158.0/24",
			"106.119.182.0/24",
			"113.207.76.0/24",
			"117.23.61.0/24",
			"118.212.233.0/24",
			"111.47.226.0/24",
			"219.153.73.0/24",
			"113.200.91.0/24",
			"203.90.247.0/24",
			"183.110.242.0/24",
			"185.254.242.0/24",
			"116.211.155.0/24",
			"116.140.35.0/24",
			"103.40.7.0/24",
			"1.255.41.0/24",
			"112.90.216.0/24",
			"1.255.100.0/24",
		},
	},
	{
		Name: "CloudFlare",
		IpCider: []string{
			"103.21.244.0/22",
			"103.22.200.0/22",
			"103.31.4.0/22",
			"104.16.0.0/13",
			"104.24.0.0/14",
			"108.162.192.0/18",
			"131.0.72.0/22",
			"141.101.64.0/18",
			"162.158.0.0/15",
			"172.64.0.0/13",
			"173.245.48.0/20",
			"188.114.96.0/20",
			"190.93.240.0/20",
			"197.234.240.0/22",
			"198.41.128.0/17",
		},
	},
	{
		Name: "云盾CDN",
		IpCider: []string{
			"27.221.64.0/24",
			"27.221.68.0/24",
			"42.236.6.128/27",
			"49.232.85.76/32",
			"58.222.57.0/28",
			"59.56.19.0/24",
			"59.56.78.0/24",
			"59.56.79.0/24",
			"60.163.162.32/27",
			"101.69.181.0/28",
			"103.95.220.0/25",
			"103.95.221.0/24",
			"103.136.251.112/28",
			"103.136.251.0/28",
			"103.219.29.64/26",
			"111.2.127.0/28",
			"111.61.59.160/27",
			"115.231.230.0/24",
			"116.136.249.0/24",
			"116.177.238.0/24",
			"117.34.43.0/24",
			"118.121.192.0/24",
			"120.53.244.232/32",
			"120.220.20.0/24",
			"122.9.54.0/24",
			"122.226.191.192/26",
			"125.44.163.0/24",
			"129.28.193.74/32",
			"153.35.236.0/28",
			"171.111.155.0/24",
			"175.6.227.128/26",
			"183.47.233.64/26",
			"183.131.145.0/28",
			"183.131.200.0/24",
			"183.134.17.0/27",
			"183.232.187.0/24",
			"223.111.172.128/28",
			"45.159.59.0/24",
			"103.100.71.0/24",
			"103.112.3.0/24",
			"117.18.111.128/25",
			"128.1.170.0/24",
			"129.227.63.0/24",
			"156.241.6.0/24",
			"164.88.96.0/24",
			"164.88.98.0/24",
			"202.181.144.128/25",
			"206.119.114.192/26",
			"206.119.110.192/26",
			"206.119.109.192/26",
			"206.119.108.192/26",
			"216.177.129.0/24",
		},
	},
	{
		Name: "腾讯云CDN",
		IpCider: []string{
			"112.29.152.0/24",
			"112.90.51.0/24",
			"113.207.39.0/24",
			"115.231.37.0/24",
			"117.169.77.0/24",
			"117.34.36.0/24",
			"119.147.227.0/24",
			"120.41.44.0/24",
			"125.39.6.0/24",
			"180.163.68.0/24",
			"182.247.229.0/24",
			"218.60.33.0/24",
			"219.146.241.0/24",
			"220.170.91.0/24",
			"221.204.182.0/24",
			"222.161.220.0/24",
			"223.87.3.0/24",
			"42.236.2.0/24",
			"58.216.25.0/24",
			"60.174.156.0/24",
			"61.184.213.0/24",
			"61.240.150.0/24",
		},
	},
}

// CheckCdnIPv4 IPv4
func CheckCdnIPv4(ipv4 string) IsCDN {
	var result IsCDN
	ip := net.ParseIP(ipv4)
	for _, v := range CdnIpList {
		for _, i := range v.IpCider {
			if _, ipNet, err2 := net.ParseCIDR(i); err2 == nil && ipNet.Contains(ip) {
				result.Is = true
				result.Name = v.Name
				return result
			}
		}
	}
	return result
}

// 常见 CDN 标志 CNAME 特征
var cdnDomains = map[string][]string{
	"Akamai": {
		"akamai-staging.net",
		"akamaized.net",
		"akamai.net",
	},
	"Amazon Cloudfront": {
		"cloudfront.net",
		"amazonaws.com",
	},
	"Ananke": {
		"anankecdn.com.br",
	},
	"Azure CDN": {
		"azureedge.net",
	},
	"Azion": {
		"azioncdn.net",
	},
	"CDN77": {
		"cdn77.net",
		"cdn77.org",
	},
	"CDNify": {
		"cdnify.io",
	},
	"CDNetworks": {
		"cdnetworks.com.gccdn.net",
	},
	"CDNsun": {
		"cdnsun.net",
	},
	"CDNvideo": {
		"cdnvideo.ru",
	},
	"Cachefly": {
		"cachefly.net",
	},
	"ChinaCache": {
		"ccgslb.com.cn",
		"chinacache.net",
		"cncssr.chinacache.net",
		"lxdns.com",
	},
	"ChinaNetCenter": {
		"wscloudcdn.com",
	},
	"CloudFlare": {
		"cdn.cloudflare.net",
	},
	"EdgeCast": {
		"edgecastcdn.net",
	},
	"Fastly": {
		"fastly.net",
	},
	"Highwinds": {
		"hwcdn.net",
	},
	"Incapsula": {
		"incapdns.net",
	},
	"Internap": {
		"internapcdn.net",
	},
	"KeyCDN": {
		"kxcdn.com",
		"awsdns",
	},
	"Leaseweb": {
		"lswcdn.net",
	},
	"Level3": {
		"fpbns.net",
		"footprint.net",
	},
	"Limelight": {
		"vo.llnwd.net",
	},
	"MaxCDN": {
		"netdna-cdn.com",
	},
	"Presscdn": {
		"presscdn.com",
	},
	"QUANTIL": {
		"mwcloudcdn.com",
		"speedcdns.com",
	},
	"Telefonica": {
		"telefonica.com",
	},
	"阿里云CDN": {
		"alikunlun.com",
		"aliyuncs.com",
		"aliyun-inc.com",
	},
	"腾讯云CDN": {
		"cdn.dnsv1.com",
		"cdntip.com",
	},
}

// CheckCNAME 检查域名 CNAME 中是否含有 CDN 关键词
func CheckCNAME(cname string) IsCDN {
	var result IsCDN
	for name, keyword := range cdnDomains {
		for _, key := range keyword {
			if strings.HasSuffix(strings.ToLower(cname), fmt.Sprintf(".%s", key)) {
				result.Is = true
				result.Name = name
				return result
			}
		}
	}
	return result
}

// 常见 CDN 响应头特征
var cdnHeaderKeywords = map[string][]string{
	"Akamai":            {"akamai", "akamai-cache-status", "x-akamai-transformed"},
	"Amazon Cloudfront": {"cloudfront", "x-amz-cf-id"},
	"Azure CDN":         {"azurecdn", "x-azure-ref"},
	"CDN77":             {"cdn77-cache"},
	"CDNetworks":        {"cdnetworks", "x-cdn"},
	"Cachefly":          {"cachefly"},
	"ChinaCache":        {"chinacache"},
	"Cloudflare":        {"cloudflare", "cf-cache-status", "cf-ray"},
	"EdgeCast":          {"edgecast", "ecs-request-id"},
	"Fastly":            {"fastly", "x-served-by"},
	"Highwinds":         {"hwcdn"},
	"Incapsula":         {"incapsula", "x-iinfo"},
	"KeyCDN":            {"keycdn"},
	"Limelight":         {"llnw", "x-llnw-cache"},
	"QUANTIL":           {"quantil", "x-quantil-cache"},
	"阿里云CDN":            {"aliyuncdn", "x-aliyun-request-id"},
	"腾讯云CDN":            {"tencentcdn", "x-tencent-cache"},
}

// CheckHeaderCDN 检查 HTTP 响应头是否含有 CDN 特征
func CheckHeaderCDN(headers map[string]string) IsCDN {
	var result IsCDN
	for cdn, keywords := range cdnHeaderKeywords {
		for _, key := range keywords {
			for headerKey, headerValue := range headers {
				if strings.Contains(strings.ToLower(headerKey), key) || strings.Contains(strings.ToLower(headerValue), key) {
					result.Is = true
					result.Name = cdn
				}
			}
		}
	}
	return result
}
