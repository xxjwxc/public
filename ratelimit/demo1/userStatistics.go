package richie

import (
"fmt"
"sort"
)

const (
	Chinese = iota
	English
)

/*
某用户剩余访问次数，例:
RemainingVisits("username")
*/
func (r *Rule) RemainingVisits(key interface{}) []int {
	arr := make([]int, 0, len(r.rules))
	for i := range r.rules {
		arr = append(arr, r.rules[i].remainingVisits(key))
	}
	return arr
}

/*
打印各细分规则下的剩余访问次数
*/
func (r *Rule) PrintRemainingVisits(key interface{}, language ...int) {
	//先确定语言，默认为中文，目前只支持中文，英文两种语言
	lan := 0
	if len(language) == 1 && language[0] == 1 {
		lan = 1
	}
	for i := range r.rules {
		if lan == 0 {
			fmt.Println(key, "在", r.rules[i].defaultExpiration, "内共允许访问", r.rules[i].numberOfAllowedAccesses, "次,剩余", r.rules[i].remainingVisits(key))
		} else {
			fmt.Println(key, "allowed", r.rules[i].numberOfAllowedAccesses, "visits within", r.rules[i].defaultExpiration, ",with", r.rules[i].remainingVisits(key), "remaining")
		}
	}
}

//获得当前所有的在线用户,注意所有用int64存储的用户会被默认认为是IP地址，会被自动转换为IP的字符串形式输出以方便查看
//如果不是本身就是以int64形式存储，而不是IP4，那么可以用ip4StringToInt64自己再转换回去
func (r *Rule) GetCurOnlineUsers() []string {
	//向切片Sli中插入没出现过的元素V，如果切片中有V，则不插入
	var insertIgnoreString = func(s []string, v string) []string {
		for _, val := range s {
			if val == v {
				return s
			}
		}
		s = append(s, v)
		return s
	}
	var users []string
	for i := range r.rules {
		f := func(k, v interface{}) bool {
			var user string
			switch k.(type) {
			case int64:
				//user = int64ToIp4String(k.(int64))
			default:
				user = fmt.Sprint(k)
			}
			users = insertIgnoreString(users, user)
			return true
		}
		r.rules[i].usedVisitorRecordsIndex.Range(f)
	}
	sort.Strings(users)
	return users
}
