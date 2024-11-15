package request

import (
	"net/http"
)

// GetHeaderArgs 获取请求头的参数
type GetHeaderArgs struct {
	Switch string
	types  string
	api    string
	Engine bool
	header map[string]string
}

// GetRequest Get请求的参数结构
type GetRequest struct {
	Timeout        int
	AllowRedirects bool
	Verify         bool
	Headers        map[string]string
	Params         map[string]string
	Stream         bool
	Engine         bool
	Data           string
	DataJson       map[string]string
	Json           map[string]any
	File           map[string][]string
}

// UrlParse URL格式化后格式
type UrlParse struct {
	Scheme   string `json:"scheme"`
	Hostname string `json:"hostname"`
	Path     string `json:"path"`
	Query    string `json:"query"`
	Port     string `json:"port"`
}

// ResponseJson 一个任意返回的Json格式
type ResponseJson map[string]interface{}

// Response 请求的返回结构
type Response struct {
	Basic struct {
		Title       string // 网页标题
		Description string // 网页描述
		Keywords    string // 网页关键字
		Favicon     string // 网页图标路径
	}
	Redirect   string                 // 重定向地址
	Url        string                 // 响应url
	StatusCode int                    // 响应状态码
	Status     string                 // 响应状态 200 ok
	Timer      float64                // 响应时长
	Headers    http.Header            // 响应头
	Body       string                 // 响应体(str)
	Charset    string                 // 检测到的编码方式
	Content    []byte                 // 响应体(byte)
	Json       map[string]interface{} // 响应的Json内容
	Length     int                    // 响应字节
	Proto      string                 // 响应协议
	ProtoMajor int                    // 响应版本号-主
	ProtoMinor int                    // 响应版本号-子
	Request    struct {
		URL     string      // 请求url
		Method  string      //请求方法
		Headers http.Header //请求头
		Body    []byte      // 请求体
	}
}

// UrlParser URL格式化结构
type UrlParser struct {
	Protocol string
	Host     string
	Port     int
	Path     string
	Params   string
}

// SiteBasic 站点基本信息结构
type SiteBasic struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Favicon     string `json:"favicon"`
}

// Fingerprint cms文件结构
type Fingerprint struct {
	Cms      string
	Method   string
	Location string
	Keyword  []string
}
