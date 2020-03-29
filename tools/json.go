package tools

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// JSONToForm tag json str to form
func JSONToForm(r *http.Request) {
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
	bodyStr := string(body)
	if len(bodyStr) > 0 {
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

// "github.com/ant0ine/go-json-rest/rest"
// func GetRequestJsonObj(r *rest.Request, v interface{}) error {

// 	//添加支持json 操作
// 	body, err := ioutil.ReadAll(r.Body)
// 	r.Body.Close()
// 	json.Unmarshal(body, &v)
// 	//-----------------------------end
// 	return err
// }

// GetJSONStr obj to json string
func GetJSONStr(obj interface{}, isFormat bool) string {
	var b []byte
	if isFormat {
		b, _ = json.MarshalIndent(obj, "", "     ")
	} else {
		b, _ = json.Marshal(obj)
	}
	return string(b)
}

// JSONDecode Json Decode
func JSONDecode(obj interface{}) string {
	return GetJSONStr(obj, false)
}

// GetJSONObj string convert to obj
func GetJSONObj(str string, out interface{}) {
	json.Unmarshal([]byte(str), out)
	return
}

// JSONEncode string convert to obj
func JSONEncode(str string, out interface{}) {
	GetJSONObj(str, out)
}
