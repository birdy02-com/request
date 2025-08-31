# request

# 使用
```
go get github.com/birdy02-com/request
```

# 介绍

`request`是一个基于 `net/http` 进行二次封装的轻量级客户端请求库，旨在为开发者提供更简洁高效的请求调用方式。它不仅简化了参数传递与请求配置，还能通过一行代码快速实现完整的请求功能。

如果你有好的建议欢迎在 [林乐天的个人博客](https://www.birdy02.com/2024/06/27/b18cf3d1-6702-4c42-b2a0-2089906d2edd) 中留言🙂。

参考文档：[https://www.birdy02.com/docs/birdy02-com-request](https://www.birdy02.com/docs/birdy02-com-request)

## 功能特点

### 1. 便捷的请求调用

- 提供更直观的 API 设计，支持传入多种参数与配置，极大简化了官方方法中繁琐的代码操作。

### 2. 丰富的响应内容

- 默认解析并提取响应中的 title、keywords、description 和 favicon 路径等关键信息。
- 支持直接获取响应体内容（官方方法通常不支持）。
- 提供 JSON、Text、Byte 等多种格式的响应数据。

### 3. 字符编码智能处理

- 集成多种方法用于检测和转换返回内容的字符编码格式，轻松避免因编码问题导致的乱码错误。

### 4. 多种请求方法支持

- 内置支持常用的 HTTP 方法：HEAD、GET、POST 和 PUT，使得开发者可以快速调用所需方法。

## 应用场景

适用于需要快速发起 HTTP 请求、处理响应内容，以及简化请求逻辑的各种场景，特别适合对高效开发和易用性有较高需求的项目。

通过 `request`，开发者可以专注于业务逻辑的实现，而无需过多关注底层请求细节，为开发工作带来极大的便利与提升。

# 主要功能
## 请求方法
1. HEAD()
2. OPTIONS()
3. GET()
4. POST()
5. PUT()

## 请求需要方法
1. GetHeader() 获取格式化请求头
2. GetPhoneHeader() 获取格式化请求头(手机)

 ## IP相关
- IsIpv4() 检查一个IP字符串是否为IPv4地址
- GetRandomIP() 生成一个随机的IP地址

## 域名方法
- RootDomain() 检查一个字符串是否为域名并获取根域名
- GetDomain() URL获取域名
- GetHostName() URL获取hostname

## URL相关  
- IsLink() 判断是否为URL
- GetUrlIpv4() 获取Url的IPv4地址
- ParseUrl() 格式化URL
- GetRootLink() 获取不带路径参数的URL
- FormatUrl() 格式化URL
- CheckCDN() 检测是否为CDN

## 其他方法
- GetAllTagA() 获取所有A标签的href值列表
- GetAllJs() 获取所有Js列表
- HttpHeaderToString() http.header 转换成 String
- HttpHeaderToMap() http.header 转换成 map

## 文件相关
- IsJsFile() 是否为Js文件
- IsCssFile() 是否为Css文件
- IsMediaFile() 是否为媒体文件

# 响应返回格式

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
		Method  string      // 请求方法
		Headers http.Header // 请求头
		Body    []byte      // 请求体
	}
	TLS *tls.ConnectionState // https请求TLS信息
}
```

# 请求使用演示
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
