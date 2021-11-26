package myhttp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

// 主要用于 爬虫

type myHttpClient struct {
	client *http.Client
}

// NewHttpClient 创建http client实例
func NewMyHttpClient() *myHttpClient {
	cookieJar, _ := cookiejar.New(nil)
	return &myHttpClient{
		client: &http.Client{
			Jar: cookieJar,
		},
	}
}

func (c *myHttpClient) GetClient() *http.Client {
	return c.client
}

// SendPost 发送POST请求
func (c *myHttpClient) SendPost(url string, params string, responseBody interface{}) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(params))
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch v := responseBody.(type) {
	case *string:
		*v = string(body)
	default:
		err = json.Unmarshal(body, responseBody)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendPost 发送POST请求
func (c *myHttpClient) SendGet(url string, params string, responseBody interface{}) error {
	req, err := http.NewRequest("GET", url, strings.NewReader(params))
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch v := responseBody.(type) {
	case *string:
		*v = string(body)
	default:
		err = json.Unmarshal(body, responseBody)
		if err != nil {
			return err
		}
	}

	return nil
}
