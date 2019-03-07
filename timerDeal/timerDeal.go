package timerDeal

import (
	"log"
	"time"
)

type TimerDeal struct {
	f_list      []func()
	f_d_timeout time.Duration
}

//增加一个回调
func (t *TimerDeal) AddOneCall(f func()) {
	t.f_list = append(t.f_list, f)
}

//设置超时时间
func (t *TimerDeal) SetCallBackTimer(d time.Duration) {
	t.f_d_timeout = d
}

func (t *TimerDeal) OnSart() {
	//time.Tick()
	go t.onTick()
}

func (t *TimerDeal) onTick() {
	ticker := time.NewTicker(t.f_d_timeout)
	for {
		time := <-ticker.C
		for _, v := range t.f_list {
			v()
		}
		log.Println("timer callback:", time.String())
	}
}
