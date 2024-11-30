# request

# 使用
```
go get github.com/birdy02-com/request
```

# 1. 简介
一款基于 net/http 二次封装的客户端请求库，支持返回更丰富的响应结构和内容；包含多种方法：
目前支持[HEAD/GET/POST]

## 更新日志
* 2024-11-30：修复处理部分url时，出现 http(s)://xxx?/的情况

# 2. 优点
* POST一键获取请求体；对比：net/http无法直接获取请求body
* 根据响应内容多层判断获取编码格式
* 默认返回目标站点的基本信息：Title、Description、Keyword、Favicon路径
* 增加请求响应时长
* 丰富默认请求头，使请求看起来像真人一样

# 3. 主要功能
1. HEAD 请求方法
2. GET 请求方法
3. POST 请求方法
4. Result 格式化响应内容
5. IsIpv4 检查一个IP字符串是否为IPv4地址
6. GetRandomIP 生成一个随机的IP地址
7. IsPrivateIP 判断IP是否为私有地址
8. ParseUrl 格式化URL
9. GetHeader 获取默认请求头
10. GetSiteBasic 提取网站的基本信息
11. GetFaviconPath 获取favicon.ico的路径

# 4. 响应返回格式

```
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
```

# 4. 请求使用演示
1. 请求参数示例
```
// GetRequest 请求的参数结构
type GetRequest struct {
	Timeout        int               // 超时时长
	AllowRedirects bool              // 是否跟随跳转
	Verify         bool              // ssl验证证书忽略
	Headers        map[string]string // 请求头
	Params         map[string]string // 请求参数
	Stream         bool
	Engine         bool
	Data           string              // post请求体
	DataJson       map[string]string   // json格式传入post请求体，会格式化成 xx=xx
	Json           map[string]any      // post请求
	File           map[string][]string // 上传的文件,格式参考 file:['文件名','内容','文件类型']
}
```
1. 简单请求
```
if r, e := request.GET(uri); e == nil {
  fmt.Println(res.Basic.Title)
}
```
2. 自定义参数请求
```
res, err := request.GET(uri, request.GetRequest{Timeout: 16})
if err != nil {
  return nil, err
}
```
3. POST请求参数
```
// 参数类型参考‘1. 请求参数示例’
res, err = request.POST(uri, request.GetRequest{Headers: Header, Data: r.Data, DataJson: r.DataJson, Json: r.Json, File: r.File})
if err == nil{
  fmt.Println(res.Basic.Title)
}
```
