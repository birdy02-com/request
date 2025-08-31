package request

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yosssi/gohtml"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var _ = os.Setenv("GODEBUG", "tlsrsakex=1")

// GetRequestInit 默认请求的参数值
func GetRequestInit() *GetRequest {
	var res GetRequest
	res.Timeout = 60
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
		for k, v := range args.Headers {
			reqArg.Headers[k] = v
		}
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
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	if args.Proxy != "" && IsLink(args.Proxy) {
		if proxyURL, err := url.Parse(args.Proxy); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	var fullURL string
	uParse, err := url.Parse(baseurl)
	if err == nil {
		// 设置Params
		params := uParse.Query()
		for k, v := range reqArg.Params {
			params.Set(k, v)
		}
		Params := params.Encode()

		if uParse.Port() != "" {
			fullURL = fmt.Sprintf("%s://%s:%s%s?%s", uParse.Scheme, uParse.Hostname(), uParse.Port(), uParse.Path, Params)
		} else {
			fullURL = fmt.Sprintf("%s://%s%s?%s", uParse.Scheme, uParse.Hostname(), uParse.Path, Params)
		}
		fullURL = strings.ReplaceAll(fullURL, " ", "%20")
		suffixToRemove := "?"
		if strings.HasSuffix(fullURL, suffixToRemove) {
			fullURL = strings.TrimSuffix(fullURL, suffixToRemove)
		}
	} else {
		fullURL = baseurl
	}
	client := http.Client{ // 创建http.Client并配置Timeout
		Timeout:   time.Duration(reqArg.Timeout) * time.Second, // 转换为time.Duration
		Transport: transport,
	}
	// 重定向策略
	if !reqArg.AllowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return reqArg, client, fullURL
}

// SetRequest 返回格式
type SetRequest struct {
	resp  *http.Response
	timer int64
	body  []byte
}

// SetGetRequest GET请求设置
func SetGetRequest(method string, client http.Client, fullURL, baseurl string, reqArg *GetRequest, args GetRequest) (*SetRequest, error) {
	method = strings.ToLower(method)
	res := SetRequest{nil, 0, nil}
	var body io.Reader
	if method == "post" || method == "put" {
		var buf bytes.Buffer
		if args.File != nil {
			w := NewWriter(&buf)
			reqArg.Headers["Content-Type"] = w.FormDataContentType()
			if args.DataJson != nil {
				for k, v := range args.DataJson {
					fileWrite, err := w.CreateFormField(k)
					if err != nil {
						return &res, err
					}
					_, err = fileWrite.Write([]byte(v))
					if err != nil {
						return &res, err
					}
				}
			}
			for Field, file := range args.File {
				if len(file) < 0 {
					return &res, errors.New("")
				}
				var err error
				var fileWrite io.Writer
				if len(file) > 2 {
					fileWrite, err = w.CreateFormFile(Field, file[0], file[2])
				} else {
					fileWrite, err = w.CreateFormFile(Field, file[0])
				}
				if err != nil {
					return &res, err
				}
				if len(file) > 1 {
					_, err = fileWrite.Write([]byte(file[1]))
				} else {
					_, err = fileWrite.Write([]byte(""))
				}
				if err != nil {
					return &res, err
				}
			}
			err := w.Close()
			if err != nil {
				return &res, err
			}
			body = &buf
			res.body = buf.Bytes()
		} else if args.Data != "" {
			body = strings.NewReader(args.Data)
			if d, e := io.ReadAll(strings.NewReader(args.Data)); e == nil {
				res.body = d
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
			res.body = buf.Bytes()
			reqArg.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		} else if args.Json != nil {
			jsonData, err := json.Marshal(args.Json)
			if err != nil {
				return &res, err
			}
			body = strings.NewReader(string(jsonData))
			res.body = jsonData
			reqArg.Headers["Content-Type"] = "application/json"
		}
	}
	var reqMethod string
	switch method {
	case "get":
		reqMethod = http.MethodGet
	case "head":
		reqMethod = http.MethodHead
	case "options":
		reqMethod = http.MethodOptions
	case "post":
		reqMethod = http.MethodPost
	case "put":
		reqMethod = http.MethodPut
	}
	req, err := http.NewRequest(reqMethod, fullURL, body)
	if err != nil {
		return &res, err
	}

	req.Header = GetHeader(&GetHeaderArgs{header: reqArg.Headers, Engine: args.Engine, api: baseurl})
	if args.Cms {
		req.AddCookie(&http.Cookie{Name: "rememberMe", Value: "me"})
	}
	timer := time.Now().UnixMicro()
	req.Host = req.Header.Get("Host")
	res.resp, err = client.Do(req)

	if args.Cms && (err != nil || res.resp == nil || res.resp.StatusCode != 200) {
		req.Header.Del("Cookie")
		timer = time.Now().UnixMicro()
		res.resp, err = client.Do(req)
		if err != nil || res.resp == nil {
			return nil, err
		}
	}
	if err != nil {
		return &res, err
	}
	res.timer = time.Now().UnixMicro() - timer
	return &res, nil
}

// Result 获取请求响应
func Result(baseurl, fullURL string, resp *http.Response, timer int64) (*Response, error) {
	var result Response
	if resp == nil {
		return &result, nil
	}
	//Header
	result.Request.URL = fullURL
	if host := resp.Request.Header.Get("Host"); host != "" {
		result.Request.URL = strings.Replace(fullURL, GetHostName(fullURL), host, -1)
	}
	result.TLS = resp.TLS
	result.Request.Method = resp.Request.Method
	result.Request.Headers = resp.Request.Header
	result.Headers = resp.Header
	result.Url = resp.Request.URL.String()
	result.Proto = resp.Proto
	result.ProtoMajor = resp.ProtoMajor
	result.ProtoMinor = resp.ProtoMinor
	result.Timer = float64(time.Now().UnixMicro()-timer) / 1e6
	result.Status = resp.Status
	result.StatusCode = resp.StatusCode
	if resp.StatusCode > 300 && resp.StatusCode < 400 {
		result.Redirect = resp.Header.Get("Location")
	}
	if resp.Request.Method == http.MethodHead {
		return &result, nil
	}

	//Body
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	var content []byte

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		// 解压 gzip
		if reader, err := gzip.NewReader(resp.Body); err == nil {
			defer func(reader *gzip.Reader) { _ = reader.Close() }(reader)
			content, _ = io.ReadAll(reader)
		}
	default:
		content, _ = io.ReadAll(resp.Body)
	}

	// 获取原始(字节)响应体

	result.Content = content
	// 获取str格式body
	body := string(content)
	body = strings.ReplaceAll(body, " ", "")
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	charsetContent, charSet := CharSetContent(content, body, contentType)
	body = html.UnescapeString(FilterString(charsetContent))
	result.Body = gohtml.Format(body)
	result.Charset = charSet
	basic := GetSiteBasic(baseurl, body)
	result.Basic.Title = basic.Title
	result.Basic.Description = basic.Description
	result.Basic.Keywords = basic.Keywords
	result.Basic.Favicon = basic.Favicon

	// 解析JSON响应到结构体
	var Json = ResponseJson{}
	_ = json.Unmarshal(content, &Json)
	if resp.Request.Body != nil {
		defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Request.Body)
		if reqBody, err := io.ReadAll(resp.Request.Body); err == nil {
			result.Request.Body = reqBody
		}
	}

	// Response
	result.Json = Json
	result.Length = len(body)
	return &result, nil
}

// HEAD 发送HTTP HEAD请求
func HEAD(baseurl string, arg ...GetRequest) (*Response, error) {
	var args GetRequest
	if len(arg) > 0 {
		args = arg[0]
	}
	reqArg, client, fullURL := GetRequestGetArg(baseurl, args)
	res, err := SetGetRequest("head", client, fullURL, baseurl, reqArg, args)
	if err != nil {
		return nil, err
	}
	result, err := Result(baseurl, fullURL, res.resp, res.timer)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// OPTIONS 发送HTTP OPTIONS请求
func OPTIONS(baseurl string, arg ...GetRequest) (*Response, error) {
	var args GetRequest
	if len(arg) > 0 {
		args = arg[0]
	}
	reqArg, client, fullURL := GetRequestGetArg(baseurl, args)
	res, err := SetGetRequest("options", client, fullURL, baseurl, reqArg, args)
	if err != nil {
		return nil, err
	}
	result, err := Result(baseurl, fullURL, res.resp, res.timer)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GET 发送HTTP GET请求
func GET(baseurl string, arg ...GetRequest) (*Response, error) {
	var args GetRequest
	if len(arg) > 0 {
		args = arg[0]
	}
	reqArg, client, fullURL := GetRequestGetArg(baseurl, args)
	res, err := SetGetRequest("get", client, fullURL, baseurl, reqArg, args)
	if err != nil {
		return nil, err
	}

	result, err := Result(baseurl, fullURL, res.resp, res.timer)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// POST 发送HTTP POST请求
func POST(baseurl string, arg ...GetRequest) (*Response, error) {
	var args GetRequest
	if len(arg) > 0 {
		args = arg[0]
	}
	reqArg, client, fullURL := GetRequestGetArg(baseurl, args)
	res, err := SetGetRequest("post", client, fullURL, baseurl, reqArg, args)
	if err != nil {
		return nil, err
	}
	result, err := Result(baseurl, fullURL, res.resp, res.timer)
	if err != nil {
		return nil, err
	}
	result.Request.Body = res.body
	return result, nil
}

// PUT 发送HTTP PUT请求
func PUT(baseurl string, arg ...GetRequest) (*Response, error) {
	var args GetRequest
	if len(arg) > 0 {
		args = arg[0]
	}
	reqArg, client, fullURL := GetRequestGetArg(baseurl, args)
	res, err := SetGetRequest("put", client, fullURL, baseurl, reqArg, args)
	if err != nil {
		return nil, err
	}
	result, err := Result(baseurl, fullURL, res.resp, res.timer)
	if err != nil {
		return nil, err
	}
	result.Request.Body = res.body
	return result, nil
}
