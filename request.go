package request

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GetRequestInit 默认请求的参数值
func GetRequestInit() *GetRequest {
	var res GetRequest
	res.Timeout = 10
	res.AllowRedirects = false
	res.Verify = false
	res.Headers = make(map[string]string)
	res.Params = make(map[string]string)
	res.Stream = true
	res.Proxy = ""
	res.Cms = false
	return &res
}

// GetRequestGetArg 获取请求参数
func GetRequestGetArg(baseurl string, args GetRequest) (*GetRequest, http.Client, string) {
	baseurl = strings.TrimSpace(baseurl)
	reqArg := GetRequestInit() // 获取请求配置
	// 处理传入参数
	// 超时时长
	if args.Timeout != 0 && args.Timeout != reqArg.Timeout {
		reqArg.Timeout = args.Timeout
	}
	// 是否跟随跳转
	if args.AllowRedirects != reqArg.AllowRedirects {
		reqArg.AllowRedirects = args.AllowRedirects
	}
	// 是否忽略证书
	if args.Verify != reqArg.Verify {
		reqArg.Verify = args.Verify
	}
	// 是否添加Headers
	if args.Headers != nil {
		reqArg.Headers = args.Headers
	}
	// 是否有Params参数
	if args.Params != nil {
		reqArg.Params = args.Params
	}
	// 是否流式传输
	if args.Stream != reqArg.Stream {
		reqArg.Stream = args.Stream
	}
	// 请求参数设置
	// 创建一个自定义的Transport，并禁用证书验证
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	if args.Proxy != "" && strings.HasPrefix(args.Proxy, "http") {
		if proxyURL, err := url.Parse(args.Proxy); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	// 设置Params
	params := url.Values{}
	for k, v := range reqArg.Params {
		params.Set(k, v)
	}
	uParse, _ := url.Parse(baseurl)
	Params := params.Encode()
	tmpParams := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(uParse.RequestURI(), uParse.Path), "?"), "/")
	if tmpParams != "" && Params != "" {
		Params = fmt.Sprintf("%s&%s", Params, tmpParams)
	} else if tmpParams != "" {
		Params = tmpParams
	}
	fullURL := fmt.Sprintf("%s://%s%s?%s", uParse.Scheme, uParse.Host, uParse.Path, Params)
	fullURL = strings.ReplaceAll(fullURL, " ", "%20")
	suffixToRemove := "?"
	if strings.HasSuffix(fullURL, suffixToRemove) {
		fullURL = strings.TrimSuffix(fullURL, suffixToRemove)
	}
	var client http.Client
	// 重定向策略
	if !reqArg.AllowRedirects {
		client = http.Client{
			Timeout: time.Duration(reqArg.Timeout) * time.Second, // 转换为time.Duration
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: transport,
		}
	} else {
		client = http.Client{ // 创建http.Client并配置Timeout
			Timeout:   time.Duration(reqArg.Timeout) * time.Second, // 转换为time.Duration
			Transport: transport,
		}
	}
	return reqArg, client, fullURL
}

// Result 获取请求响应
func Result(baseurl, fullURL string, resp *http.Response) (*Response, error) {
	var result Response
	defer func(Body io.ReadCloser) {
		if err1 := Body.Close(); err1 != nil {
		}
	}(resp.Body)
	// 获取原始(字节)响应体
	content, _ := io.ReadAll(resp.Body)
	result.Content = content
	// 获取str格式body
	body := string(content)
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	charsetContent, charSet := CharSetContent(content, body, contentType)
	body = html.UnescapeString(charsetContent)
	result.Body = body
	result.Charset = charSet
	basic := GetSiteBasic(baseurl, body)
	result.Basic.Title = basic.Title
	result.Basic.Description = basic.Description
	result.Basic.Keywords = basic.Keywords
	result.Basic.Favicon = basic.Favicon
	// 解析JSON响应到结构体
	var Json = ResponseJson{}
	_ = json.Unmarshal(content, &Json)
	// Request 参数
	result.Request.URL = fullURL
	result.Request.Method = resp.Request.Method
	result.Request.Headers = resp.Request.Header
	if resp.Request.Method == "POST" && resp.Request.Body != nil {
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
			}
		}(resp.Request.Body)
		if reqBody, err := io.ReadAll(resp.Request.Body); err == nil {
			result.Request.Body = reqBody
		}
	}

	// Response
	result.Status = resp.Status
	result.StatusCode = resp.StatusCode
	if resp.StatusCode > 300 && resp.StatusCode < 400 {
		result.Redirect = resp.Header.Get("Location")
	}
	result.Headers = resp.Header
	result.Url = resp.Request.URL.String()

	result.Json = Json
	result.Length = len(body)
	result.Proto = resp.Proto
	result.ProtoMajor = resp.ProtoMajor
	result.ProtoMinor = resp.ProtoMinor
	return &result, nil
}

// HEAD 发送HTTP HEAD请求
func HEAD(baseurl string, arg ...GetRequest) (*Response, error) {
	var args GetRequest
	if len(arg) > 0 {
		args = arg[0]
	}
	reqArg, client, fullURL := GetRequestGetArg(baseurl, args)
	req, err := http.NewRequest(http.MethodHead, fullURL, nil)
	if err != nil {
		return &Response{}, err
	}
	req.Header = GetHeader(&GetHeaderArgs{header: reqArg.Headers, Engine: args.Engine, api: baseurl})
	timer := time.Now().UnixMicro()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	timer = time.Now().UnixMicro() - timer
	result, err := Result(baseurl, fullURL, resp)
	if err != nil {
		return nil, err
	}
	result.Timer = float64(time.Now().UnixMicro()-timer) / 1e6
	return result, nil
}

// OPTIONS 发送HTTP OPTIONS请求
func OPTIONS(baseurl string, arg ...GetRequest) (*Response, error) {
	var args GetRequest
	if len(arg) > 0 {
		args = arg[0]
	}
	reqArg, client, fullURL := GetRequestGetArg(baseurl, args)
	req, err := http.NewRequest(http.MethodHead, fullURL, nil)
	if err != nil {
		return &Response{}, err
	}
	req.Header = GetHeader(&GetHeaderArgs{header: reqArg.Headers, Engine: args.Engine, api: baseurl})
	timer := time.Now().UnixMicro()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	timer = time.Now().UnixMicro() - timer
	result, err := Result(baseurl, fullURL, resp)
	if err != nil {
		return nil, err
	}
	result.Timer = float64(time.Now().UnixMicro()-timer) / 1e6
	return result, nil
}

// GET 发送HTTP GET请求
func GET(baseurl string, arg ...GetRequest) (*Response, error) {
	var args GetRequest
	if len(arg) > 0 {
		args = arg[0]
	}
	reqArg, client, fullURL := GetRequestGetArg(baseurl, args)
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return &Response{}, err
	}
	req.Header = GetHeader(&GetHeaderArgs{header: reqArg.Headers, Engine: args.Engine, api: baseurl})
	if args.Cms {
		req.AddCookie(&http.Cookie{Name: "rememberMe", Value: "me"})
	}
	timer := time.Now().UnixMicro()
	resp, err := client.Do(req)
	if args.Cms && (err != nil || resp != nil && resp.StatusCode != 200) {
		req.Header.Del("Cookie")
		timer = time.Now().UnixMicro()
		resp, err = client.Do(req)
		if err != nil || resp == nil {
			return nil, err
		}
	}
	result, err := Result(baseurl, fullURL, resp)
	if err != nil {
		return nil, err
	}
	result.Timer = float64(time.Now().UnixMicro()-timer) / 1e6
	return result, nil
}

// POST 发送HTTP POST请求
func POST(baseurl string, arg ...GetRequest) (*Response, error) {
	var args GetRequest
	if len(arg) > 0 {
		args = arg[0]
	}
	reqArg, client, fullURL := GetRequestGetArg(baseurl, args)
	var body io.Reader
	var buf bytes.Buffer
	var bodyByte = make([]byte, 0)
	if args.File != nil {
		w := NewWriter(&buf)
		reqArg.Headers["Content-Type"] = w.FormDataContentType()
		if args.DataJson != nil {
			for k, v := range args.DataJson {
				fileWrite, err := w.CreateFormField(k)
				if err != nil {
					return nil, err
				}
				_, err = fileWrite.Write([]byte(v))
				if err != nil {
					return nil, err
				}
			}
		}
		for Field, file := range args.File {
			if len(file) < 0 {
				return nil, errors.New("")
			}
			var err error
			var fileWrite io.Writer
			if len(file) > 2 {
				fileWrite, err = w.CreateFormFile(Field, file[0], file[2])
			} else {
				fileWrite, err = w.CreateFormFile(Field, file[0])
			}
			if err != nil {
				return nil, err
			}
			if len(file) > 1 {
				_, err = fileWrite.Write([]byte(file[1]))
			} else {
				_, err = fileWrite.Write([]byte(""))
			}
			if err != nil {
				return nil, err
			}
		}
		err := w.Close()
		if err != nil {
			return nil, err
		}
		body = &buf
		bodyByte = buf.Bytes()
	} else if args.Data != "" {
		body = strings.NewReader(args.Data)
		if d, e := io.ReadAll(strings.NewReader(args.Data)); e == nil {
			bodyByte = d
		}
		reqArg.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	} else if args.DataJson != nil {
		isFirst := true
		for key, value := range args.DataJson {
			if !isFirst {
				buf.WriteString("&")
			}
			isFirst = false
			buf.WriteString(url.QueryEscape(key))
			buf.WriteString("=")
			buf.WriteString(url.QueryEscape(value))
		}
		body = &buf
		bodyByte = buf.Bytes()
		reqArg.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	} else if args.Json != nil {
		jsonData, err := json.Marshal(args.Json)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(string(jsonData))
		bodyByte = jsonData
		reqArg.Headers["Content-Type"] = "application/json"
	}
	req, err := http.NewRequest(http.MethodPost, fullURL, body)
	if err != nil {
		return &Response{}, err
	}
	req.Header = GetHeader(&GetHeaderArgs{header: reqArg.Headers, Engine: args.Engine, api: baseurl})
	timer := time.Now().UnixMicro()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	result, err := Result(baseurl, fullURL, resp)
	if err != nil {
		return nil, err
	}
	result.Timer = float64(time.Now().UnixMicro()-timer) / 1e6
	result.Request.Body = bodyByte
	return result, nil
}

// PUT 发送HTTP PUT请求
func PUT(baseurl string, arg ...GetRequest) (*Response, error) {
	var args GetRequest
	if len(arg) > 0 {
		args = arg[0]
	}
	reqArg, client, fullURL := GetRequestGetArg(baseurl, args)
	var body io.Reader
	var buf bytes.Buffer
	var bodyByte = make([]byte, 0)
	if args.File != nil {
		w := NewWriter(&buf)
		reqArg.Headers["Content-Type"] = w.FormDataContentType()
		if args.DataJson != nil {
			for k, v := range args.DataJson {
				fileWrite, err := w.CreateFormField(k)
				if err != nil {
					return nil, err
				}
				_, err = fileWrite.Write([]byte(v))
				if err != nil {
					return nil, err
				}
			}
		}
		for Field, file := range args.File {
			if len(file) < 0 {
				return nil, errors.New("")
			}
			var err error
			var fileWrite io.Writer
			if len(file) > 2 {
				fileWrite, err = w.CreateFormFile(Field, file[0], file[2])
			} else {
				fileWrite, err = w.CreateFormFile(Field, file[0])
			}
			if err != nil {
				return nil, err
			}
			if len(file) > 1 {
				_, err = fileWrite.Write([]byte(file[1]))
			} else {
				_, err = fileWrite.Write([]byte(""))
			}
			if err != nil {
				return nil, err
			}
		}
		err := w.Close()
		if err != nil {
			return nil, err
		}
		body = &buf
		bodyByte = buf.Bytes()
	} else if args.Data != "" {
		body = strings.NewReader(args.Data)
		if d, e := io.ReadAll(strings.NewReader(args.Data)); e == nil {
			bodyByte = d
		}
		reqArg.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	} else if args.DataJson != nil {
		isFirst := true
		for key, value := range args.DataJson {
			if !isFirst {
				buf.WriteString("&")
			}
			isFirst = false
			buf.WriteString(url.QueryEscape(key))
			buf.WriteString("=")
			buf.WriteString(url.QueryEscape(value))
		}
		body = &buf
		bodyByte = buf.Bytes()
		reqArg.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	} else if args.Json != nil {
		jsonData, err := json.Marshal(args.Json)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(string(jsonData))
		bodyByte = jsonData
		reqArg.Headers["Content-Type"] = "application/json"
	}
	req, err := http.NewRequest(http.MethodPut, fullURL, body)
	if err != nil {
		return &Response{}, err
	}
	req.Header = GetHeader(&GetHeaderArgs{header: reqArg.Headers, Engine: args.Engine, api: baseurl})
	timer := time.Now().UnixMicro()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	result, err := Result(baseurl, fullURL, resp)
	if err != nil {
		return nil, err
	}
	result.Timer = float64(time.Now().UnixMicro()-timer) / 1e6
	result.Request.Body = bodyByte
	return result, nil
}
