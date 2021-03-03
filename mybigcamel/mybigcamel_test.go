package mybigcamel

import (
	"fmt"
	"strings"
	"testing"
)

func Test_cache(t *testing.T) {
	SS := "OauthIDAPI"

	tmp0 := UnMarshal(SS)
	fmt.Println(tmp0)
	tmp1 := Marshal(tmp0)
	fmt.Println(tmp1)

	if SS != tmp1 {
		fmt.Println("false.")
	}

	fmt.Println(CapLowercase("IDAPIID"))
	fmt.Println(CapSmallcase("IDAPIID"))
}

func CapLowercase(name string) string {
	list := strings.Split(UnMarshal(name), "_")
	if len(list) == 0 {
		return ""
	}

	return list[0] + name[len(list[0]):]
}

func CapSmallcase(name string) string {
	list := strings.Split(UnSmallMarshal(name), "_")
	if len(list) == 0 {
		return ""
	}

	return list[0] + name[len(list[0]):]
}
