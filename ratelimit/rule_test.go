package ratelimit

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	var Visits int //因并发问题num比实际数量稍小
	fmt.Println("\r\n测试1,性能测试，预计耗时1分钟，请耐心等待:")
	//步骤一：初始化
	r := NewRule()
	//步骤二：增加一条或者多条规则组成复合规则，规则必须至少包含一条规则
	//此处对于性能测试，为方便准确计数，只需要添加一条规则
	err := r.AddRule(time.Second*10, 1000) //每10秒只允许访问1000次
	if err != nil {
		panic(err)
	}
	/*
		r.AddRule(time.Second*10, 10)   //每10秒只允许访问10次
		r.AddRule(time.Minute*30, 1000) //每30分钟只允许访问1000次
		r.AddRule(time.Hour*24, 5000)   //每天只允许访问500次
	*/
	//步骤三(可选):从本地磁盘加载历史访问数据
	//r.LoadingAndAutoSaveToDisc("test1", time.Second*10) //设置10秒备份一次(不填写则默认60秒备份一次)，备份到程序当前文件夹下，文件名为test1.ratelimit
	log.Println("性能测试正式开始")
	//步骤四：调用函数判断某用户是否允许访问
	/*
	   allow:= r.AllowVisit(user)
	*/
	//构建若干个用户，模拟用户访问
	var users = make(map[string]bool)
	for i := 1; i < 1000; i++ {
		users["user_"+strconv.Itoa(i)] = true
	}
	begin := time.Now()
	//模拟多个协程访问
	chanNum := 200
	var wg sync.WaitGroup
	wg.Add(chanNum)
	for i := 0; i < chanNum; i++ {
		go func(i int, wg *sync.WaitGroup) {
			for ii := 0; ii < 5000; ii++ {
				for user := range users {
					for {
						Visits++
						b, err := r.AllowVisit(user)
						if err != nil {
							panic(err)
						}
						if !b {
							break
						}
					}
				}
			}
			wg.Done()
		}(i, &wg)
	}
	//所有线程结束，完工
	wg.Wait()
	t1 := int(time.Now().Sub(begin).Seconds())
	log.Println("性能测试完成:共计访问", Visits, "次,", "耗时", t1, "秒,即每秒约完成", Visits/t1, "次操作")
	//步骤五(可选):程序退出前主动手动存盘
	//err := r.SaveToDiscOnce() //在自动备份的同时，还支持手动备份，一般在程序要退出时调用此函数
	//if err == nil {
	//	log.Println("完成手动数据备份")
	//} else {
	//	log.Println(err)
	//}
}

func Test2(t *testing.T) {
	fmt.Println("\r\n测试2，模拟用户访问并打印:")
	//步骤一：初始化
	r := NewRule()
	//步骤二：增加一条或者多条规则组成复合规则，规则必须至少包含一条规则
	r.AddRule(time.Second*10, 0)  //每10秒只允许访问5次
	r.AddRule(time.Minute*30, 50) //每30分钟只允许访问50次
	r.AddRule(time.Hour*24, 500)  //每天只允许访问500次
	//步骤三：调用函数判断某用户是否允许访问
	/*
	   allow:= r.AllowVisit(user)
	*/
	//构建若干个用户，模拟用户访问
	users := []string{"andyyu", "tony", "chery"}
	for _, user := range users {
		fmt.Println("\r\n开始模拟以下用户访问:", user)
		for {
			b, err := r.AllowVisit(user)
			if err != nil {
				panic(err)
			}
			if b {
				log.Println(user, "访问1次,剩余:", r.RemainingVisits(user))
			} else {
				log.Println(user, "访问过多,稍后再试")
				break
			}
			time.Sleep(time.Second * 1)
		}
	}
	//打印所有用户访问数据情况
	fmt.Println("开始打印所有用户在相关时间段内详细的剩余访问次数情况:")
	for _, user := range users {
		fmt.Println(user)
		fmt.Println("     概述:", r.RemainingVisits(user))
		fmt.Println("     具体:")
		r.PrintRemainingVisits(user)
		fmt.Println("")
	}
	/*
		在实际的平台运行过程中，往往会因为各种原因，某个客户的访问量过大，被系统临时禁止访问，这时候
		这个客户就可能会投诉之类的，根据运营的实际需要，就需要手动清除掉某用户的访问记录，让其可以再继续访问。
		对于函数ManualEmptyVisitorRecordsOf(),一般需要自行通过合理的方式,比如自行封装一个HTTP服务来间接调用
	*/
	log.Println("开始测试手动清楚某用户访问记录.")
	log.Println("chery清空访问记录前,剩余:", r.RemainingVisits("chery"))
	r.ManualEmptyVisitorRecordsOf("chery")
	log.Println("chery清空访问记录后,剩余:", r.RemainingVisits("chery"))

}
