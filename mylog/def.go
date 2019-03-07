package mylog

import (
	"time"
)

/*
 日志详细信息
*/
type LogInfo struct {
	Service    string      `json:"service"` //服务名
	Group      string      `json:"group"`   //服务组
	Type       string      `json:"type"`    //日志的类型()
	Action     string      `json:"action"`  //动作
	Path       string      `json:"path"`    //路径,地址,深度
	Ip         string      `json:"ip"`      //ip地址
	Topic      string      `json:"topic"`
	Bundle     string      `json:"bundle"`
	Pid        string      `json:"pid"`
	Data       interface{} `json:"data"`
	Creat_time time.Time   `json:"created"`
}

/*
向es发送数据结构信息
*/
type EsLogInfo struct {
	Info     LogInfo `json:"loginfo"`  //信息
	Es_index string  `json:"es_index"` //索引
	Es_type  string  `json:"es_type"`  //类型
	Es_id    string  `json:"es_id"`    //id
}

const (
	Http_log_index = "http_log"
)

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 5
	},
	"mappings":{
		"` + Http_log_index + `":{
			"properties":{
				"service":{
					"type":"keyword"
				},
				"group":{
					"type":"keyword"
				},
				"type":{
					"type":"text"
				},
				"action":{
					"type":"keyword"
				},
				"path":{
					"type":"text"
				},
				"ip":{
					"type":"text"
				},
				"topic":{
					"type":"keyword"
				},
				"bundle":{
					"type":"keyword"
				},
				"pid":{
					"type":"keyword"
				},
				"data":{
					"type":"text"
				},
				"created":{
					"type":"date"
				}
			}
		}
	}
}`
