package request

import (
	"bytes"
	"fmt"
	"github.com/yinheli/mahonia"
	"golang.org/x/net/html/charset"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ---------- User-Agent ----------

// GetHeader 获取请求头
func GetHeader(args *GetHeaderArgs) http.Header {
	// 预设的User-Agent列表
	userAgents := []string{
		// Chrome 系列
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.112 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.207 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.76 Safari/537.36",

		// Edge 系列
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.112 Safari/537.36 Edg/125.0.2535.67",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.207 Safari/537.36 Edg/124.0.2478.97",

		// Firefox 系列
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:126.0) Gecko/20100101 Firefox/126.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.4; rv:126.0) Gecko/20100101 Firefox/126.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:126.0) Gecko/20100101 Firefox/126.0",

		// Safari (Mac)
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",

		// Opera 系列
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.112 Safari/537.36 OPR/110.0.4985.67",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.207 Safari/537.36 OPR/109.0.4971.80",

		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.112 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.76 Safari/537.36",
		"Mozilla/5.0 (Windows NT 11.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.6480.42 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/125.0.2535.67 Chrome/125.0.6422.112 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.112 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:126.0) Gecko/20100101 Firefox/126.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.5; rv:126.0) Gecko/20100101 Firefox/126.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) OPR/110.0.4985.67 Chrome/125.0.6422.112 Safari/537.36",
	}
	// 初始化HTTP头
	header := http.Header{}
	rand.Seed(time.Now().UnixNano())
	if args.Engine {
		header.Set("User-Agent", "Baiduspider+(+https://www.baidu.com/search/spider.htm);google|baiduspider|baidu|spider|sogou|bing|yahoo|soso|sosospider|360spider|youdao|jikeSpider;)")
	} else {
		header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	}
	header.Set("Connection", "keep-alive")
	header.Set("Cache-Control", "max-age=0")
	header.Set("Upgrade-Insecure-Requests", "1")
	header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	//header.Set("Accept-Encoding", "gzip, deflate, br")
	header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	header.Set("Sec-Ch-Ua", `"Chromium";v="111", "Not(A:Brand";v="8"`)
	header.Set("Sec-Ch-Ua-Mobile", "?0")
	header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	header.Set("Sec-Fetch-Site", `none`)
	header.Set("Sec-Fetch-Mode", `navigate`)
	header.Set("Sec-Fetch-User", `?1`)
	header.Set("Sec-Fetch-Dest", `document`)

	// 根据参数更新HTTP头
	if args.Switch == "Bing" {
		header.Set("Host", "cn.bing.com")
		header.Set("Referer", "https://www.bing.com")
	}
	if args.api != "" {
		if parse, err := url.Parse(args.api); err == nil {
			header.Set("Host", parse.Host) // 提取Host和Referer
		}
		header.Set("Referer", args.api)
	}
	// 传入的header参数
	for k, v := range args.header {
		header.Set(k, v)
	}
	return header
}

// GetPhoneHeader 获取Phone请求头
func GetPhoneHeader() string {
	// 预设的User-Agent列表
	userAgents := []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25",
		"Mozilla/5.0 (iPad; CPU OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25",
		"Mozilla/5.0 (Linux; Android 4.4.2; Nexus 4 Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.114 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.199 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 14; SAMSUNG SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/22.0 Chrome/120.0.6099.199 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 13; Mi 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.199 Mobile Safari/537.36 MIUIBrowser/14.0.6",
		"Mozilla/5.0 (Linux; Android 12; HUAWEI P50) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.199 Mobile Safari/537.36 HuaweiBrowser/12.0.3.301",
	}
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
}

// ---------- Decode HTML ----------

// CharSetContent 编码后的text
func CharSetContent(content []byte, body string, contentType string) (string, string) {
	var htmlEncode string
	// BOM检测编码
	if DetectBOM(content) {
		return Convert(body, "utf-8", "utf-8"), "utf-8"
	}
	// 响应头检测编码
	contentType = strings.ToLower(contentType)
	if contentType != "" {
		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
			htmlEncode = "gbk"
		} else if strings.Contains(contentType, "big5") {
			htmlEncode = "big5"
		} else if strings.Contains(contentType, "utf-8") {
			htmlEncode = "utf-8"
		}
	}
	if htmlEncode != "" {
		return Convert(body, htmlEncode, "utf-8"), htmlEncode
	}
	// 匹配正文中的编码
	match := regexp.MustCompile(`(?is)<meta[^>]*charset\s*=["']?\s*([A-Za-z0-9\-]+)`).FindStringSubmatch(body)
	if len(match) > 1 {
		contentType = strings.ToLower(match[1])
		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
			htmlEncode = "gbk"
		} else if strings.Contains(contentType, "big5") {
			htmlEncode = "big5"
		} else if strings.Contains(contentType, "utf-8") {
			htmlEncode = "utf-8"
		}
	}
	if htmlEncode != "" {
		return Convert(body, htmlEncode, "utf-8"), htmlEncode
	}
	// 自动检测编码
	_, contentType, _ = charset.DetermineEncoding(content, "")
	if contentType != "" {
		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
			htmlEncode = "gbk"
		} else if strings.Contains(contentType, "big5") {
			htmlEncode = "big5"
		} else if strings.Contains(contentType, "utf-8") {
			htmlEncode = "utf-8"
		}
	}
	if htmlEncode != "" {
		return Convert(body, htmlEncode, "utf-8"), htmlEncode
	}
	// 默认返回utf-8
	return Convert(body, "utf-8", "utf-8"), htmlEncode
}

// DetectBOM 检测是否为UTF-8 BOM
func DetectBOM(data []byte) bool {
	return bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF})
}

// Convert 编码HTML
func Convert(src string, srcCode string, tagCode string) string {
	if srcCode == tagCode {
		return src
	}
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// FilterString 过滤非法字符串
func FilterString(input string) string {
	var output []rune
	for _, r := range input {
		if (r >= 0xD800 && r <= 0xDFFF) || (r >= 0xE000 && r <= 0xF8FF) || r == 0xFFFE || r == 0xFFFF {
			continue
		}
		output = append(output, r)
	}
	return string(output)
}

// ---------- Basic Info ----------

// GetSiteBasic 提取网站基本信息
func GetSiteBasic(baseurl, text string) *SiteBasic {
	var result SiteBasic
	reTitle := regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)

	// 查找所有匹配项
	matchesTitle := reTitle.FindStringSubmatch(text)
	if len(matchesTitle) > 0 {
		title := matchesTitle[1]
		result.Title = title
	}

	reDes := regexp.MustCompile(`<meta\s+name=["']description["']\s+content=["'](.*?)["']\s*/?>`)
	// 查找所有匹配项
	matchesDes := reDes.FindStringSubmatch(text)
	if len(matchesDes) > 0 {
		desc := matchesDes[1]
		result.Description = desc
	}

	reKeywords := regexp.MustCompile(`<meta\s+name=["']keywords["']\s+content=["'](.*?)["']\s*/?>`)
	// 查找所有匹配项
	matchesKeywords := reKeywords.FindStringSubmatch(text)
	if len(matchesKeywords) > 0 {
		keywords := matchesKeywords[1]
		result.Keywords = keywords
	}

	favicon := GetFaviconPath(baseurl, text)
	result.Favicon = favicon
	return &result
}

// GetFaviconPath 获取favicon.ico的路径
func GetFaviconPath(uri, body string) string {
	//regFav := regexp.MustCompile(`href="(.*?favicon....)"`)
	regFav := regexp.MustCompile(`rel="icon" href="(.*?favicon[^"]*)">`)
	matchFav := regFav.FindAllStringSubmatch(body, -1)
	if len(matchFav) < 1 {
		regFav = regexp.MustCompile(`type="image/x-icon" href="(.*?favicon[^"]*)">`)
		matchFav = regFav.FindAllStringSubmatch(body, -1)
	}
	var faviconPath string
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	uri = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	if len(matchFav) > 0 {
		fav := matchFav[0][1]
		if fav[:2] == "//" {
			faviconPath = fmt.Sprintf("http:%s", fav)
		} else {
			if fav[:4] == "http" {
				faviconPath = fav
			} else {
				faviconPath = fmt.Sprintf("%s/%s", uri, strings.TrimPrefix(fav, "/"))
			}
		}
	} else {
		faviconPath = fmt.Sprintf("%s/favicon.ico", uri)
	}
	return faviconPath
}
