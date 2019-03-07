package tools

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

func JsonToForm(r *http.Request) {
	//添加支持json 操作
	r.ParseForm()
	if len(r.Form) == 1 { //可能是json 支持json
		for key, value := range r.Form {
			if len(value[0]) == 0 {
				delete(r.Form, key)
				var m map[string]string
				if err := json.Unmarshal([]byte(key), &m); err == nil {
					for k, v := range m {
						r.Form[k] = []string{v}
					}
				}
			}
		}
	}

	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	if len(body_str) > 0 {
		var m map[string]string
		if err := json.Unmarshal(body, &m); err == nil {
			for k, v := range m {
				r.Form[k] = []string{v}
			}
		}
	}
	//-----------------------------end
	return
}

func GetRequestJsonObj(r *rest.Request, v interface{}) error {

	//添加支持json 操作
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	json.Unmarshal(body, &v)
	//-----------------------------end
	return err
}

func GetJsonStr(obj interface{}) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

func JsonDecode(obj interface{}) string {
	return GetJsonStr(obj)
}

func GetJsonObj(str string, out interface{}) {
	json.Unmarshal([]byte(str), out)
	return
}

func JsonEncode(str string, out interface{}) {
	GetJsonObj(str, out)
}
