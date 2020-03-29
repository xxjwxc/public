package server

import (
	"fmt"
	"os/exec"
	"strings"
)

type ServiceTools struct {
	//i IServiceTools
}

func (s *ServiceTools) IsStart(name string) (st int, err error) {
	f, _ := exec.Command("service", name, "status").Output()

	st = NOTFIND
	str := string(f)
	a := strings.Split(str, "\n")
	for _, v := range a {
		if strings.Index(v, "Active:") > 0 {
			fmt.Println("====info===:", v)
			if strings.Index(v, "inactive") > 0 { //不活动的
				st = Stopped
			} else if strings.Index(v, "activating") > 0 { //活动的
				st = Running
			}
			break
		}
	}

	return
}
