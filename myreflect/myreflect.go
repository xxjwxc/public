package myreflect

import (
	"reflect"
	"strings"
)

// FindTag find struct of tag string.查找struct 的tag信息
func FindTag(obj interface{}, field, tag string) string {
	dataStructType := reflect.Indirect(reflect.ValueOf(obj)).Type()
	for i := 0; i < dataStructType.NumField(); i++ {
		fd := dataStructType.Field(i)
		if fd.Name == field {
			bb := fd.Tag
			sqlTag := bb.Get(tag)

			if sqlTag == "-" || bb == "-" {
				return ""
			}

			sqlTags := strings.Split(sqlTag, ",")
			sqlFieldName := fd.Name // default
			if len(sqlTags[0]) > 0 {
				sqlFieldName = sqlTags[0]
			}
			return sqlFieldName
		}
	}

	return ""
}
