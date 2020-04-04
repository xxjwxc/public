package myswagger

import "strings"

var version string = "2.0"
var host string = "localhost"
var basePath string = "/v1"
var schemes []string = []string{"http", "https"}
var reqCtxType []string = []string{"application/json", "application/xml"}
var respCtxType []string = []string{"application/json", "application/xml"}
var info Info = Info{
	Description: "swagger default desc",
	Version:     "1.0.0",
	Title:       "Swagger Petstore",
}
var externalDocs ExternalDocs = ExternalDocs{
	Description: "Find out more about Swagger",
	URL:         "https://github.com/xxjwxc/public",
}

// SetVersion 设置版本号
func SetVersion(v string) {
	version = v
}

// SetHost 设置host
func SetHost(h string) {
	h = strings.TrimPrefix(h, "http://")
	h = strings.TrimPrefix(h, "https://")
	host = h
}

// GetHost 获取host
func GetHost() string {
	return schemes[0] + "://" + host
}

// SetBasePath set basePath
func SetBasePath(b string) {
	if !strings.HasPrefix(b, "/") {
		b = "/" + b
	}
	basePath = b
}

// SetSchemes 设置 http头
func SetSchemes(isHTTP, isHTTPS bool) {
	schemes = []string{}
	if isHTTP {
		schemes = append(schemes, "http")
	}
	if isHTTPS {
		schemes = append(schemes, "https")
	}
}

// SetReqCtxType 设置请求数据传输方式
func SetReqCtxType(isJSON, isXML bool) {
	reqCtxType = []string{}
	if isJSON {
		reqCtxType = append(schemes, "application/json")
	}
	if isXML {
		reqCtxType = append(schemes, "application/xml")
	}
}

// SetRespCtxType 设置响应(返回)求数据传输方式
func SetRespCtxType(isJSON, isXML bool) {
	respCtxType = []string{}
	if isJSON {
		respCtxType = append(schemes, "application/json")
	}
	if isXML {
		respCtxType = append(schemes, "application/xml")
	}
}

// SetInfo 设置信息
func SetInfo(i Info) {
	info = i
}

// SetExternalDocs 设置外部doc链接
func SetExternalDocs(e ExternalDocs) {
	externalDocs = e
}
