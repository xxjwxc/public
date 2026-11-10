package myhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/xxjwxc/public/mylog"
)

// DefaultTimeout 默认 http 请求超时时间。为 0 表示不限制超时(与旧版本行为一致)。
// 调用方可按需修改,如 myhttp.DefaultTimeout = 10*time.Second
var DefaultTimeout = 30 * time.Second

// newHTTPClient 创建一个带默认超时的 http 客户端。
// 底层复用 http.DefaultTransport,因此连接池在各次请求间共享。
func newHTTPClient() *http.Client {
	return &http.Client{Timeout: DefaultTimeout}
}

// OnPostJSON 发送修改密码
func OnPostJSON(url, jsonstr string) ([]byte, error) {
	//解析这个 URL 并确保解析没有出错。
	body := bytes.NewBuffer([]byte(jsonstr))
	resp, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body1, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil, err
	}

	return body1, nil
}

// OnGetJSON 发送get 请求
func OnGetJSON(url, params string) string {
	//解析这个 URL 并确保解析没有出错。
	var urls = url
	if len(params) > 0 {
		urls += "?" + params
	}
	resp, err := http.Get(urls)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body1, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		mylog.Error(url, "--1--", err1)
		return ""
	}

	return string(body1)
}

// OnGetJSON2 发送get 请求
func OnGetJSON2(url, params string) string {
	//解析这个 URL 并确保解析没有出错。
	var urls = url
	if len(params) > 0 {
		urls += "?" + params
	}
	resp, err := http.Get(urls)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body1, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		mylog.Error(url, err1)
		time.Sleep(3 * time.Second)
		return OnGetJSON(url, params) // 重试
	}

	return string(body1)
}

// SendGet 发送get 请求 返回对象
func SendGet(url, params string, obj interface{}) bool {
	//解析这个 URL 并确保解析没有出错。
	var urls = url
	if len(params) > 0 {
		urls += "?" + params
	}
	resp, err := http.Get(urls)
	if err != nil {
		mylog.Error(err)
		return false
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		mylog.Error(err)
		return false
	}
	//log.Println((string(body)))
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		mylog.Error(err)
		return false
	}

	return true
}

// SendGetEx 发送GET请求
func SendGetEx(url string, reponse interface{}) bool {
	resp, e := http.Get(url)
	if e != nil {
		mylog.Error(e)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Error(err)
		return false
	}
	//mylog.Debug(string(body))
	err = json.Unmarshal(body, &reponse)
	if err != nil {
		mylog.Error(err)
		return false
	}

	return true
}

// OnPostForm form 方式发送post请求
func OnPostForm(url string, data url.Values) (body []byte) {
	resp, err := http.PostForm(url, data)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

// SendPost 发送POST请求
func SendPost(requestBody interface{}, responseBody interface{}, url string) error {
	postData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(postData))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	//	req.Header.Add("Authorization", authorization)
	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// mylog.Debug(string(body))

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSON  像指定client 发送json 包
// msg message.MessageBody
func WriteJSON(w http.ResponseWriter, msg interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	js, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(js))
}

// PostHeader post by head
func PostHeader(url, jsonstr string, header http.Header) ([]byte, error) {
	postData := bytes.NewBuffer([]byte(jsonstr))
	req, err := http.NewRequest("POST", url, postData)
	if err != nil {
		return nil, err
	}
	req.Header = header
	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetHeader get by head
func GetHeader(url string, header http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header
	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
