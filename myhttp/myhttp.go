package myhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/xxjwxc/public/mylog"
)

//OnPostJSON 发送修改密码
func OnPostJSON(url, jsonstr string) []byte {
	//解析这个 URL 并确保解析没有出错。
	body := bytes.NewBuffer([]byte(jsonstr))
	resp, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		return []byte("")
	}
	defer resp.Body.Close()
	body1, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		mylog.Error(err1)
		return []byte("")
	}

	return body1
}

//OnGetJSON 发送get 请求
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
		mylog.Error(err1)
		return ""
	}

	return string(body1)
}

//SendGet 发送get 请求 返回对象
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
	body, err := ioutil.ReadAll(resp.Body)
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

//SendGetEx 发送GET请求
func SendGetEx(url string, reponse interface{}) bool {
	resp, e := http.Get(url)
	if e != nil {
		mylog.Error(e)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Error(e)
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

//OnPostForm form 方式发送post请求
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

//SendPost 发送POST请求
func SendPost(requestBody interface{}, responseBody interface{}, url string) bool {
	postData, err := json.Marshal(requestBody)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(postData))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	//	req.Header.Add("Authorization", authorization)
	resp, e := client.Do(req)
	if e != nil {
		mylog.Error(e)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Error(e)
		return false
	}
	//	result := string(body)
	//mylog.Debug(string(body))

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		mylog.Error(err)
		return false
	}

	return true
}

//WriteJSON  像指定client 发送json 包
//msg message.MessageBody
func WriteJSON(w http.ResponseWriter, msg interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	js, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(js))
}
