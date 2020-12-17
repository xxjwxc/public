package mynotify

import (
	"time"
)

type notify struct {
	tag int
	cc  chan int // -2:超时退出 -1:已关闭 0:默认状态 1:正在运行
	tm  time.Duration
}

// New new signal
func New(tm time.Duration) *notify {
	if tm == 0 { // 默认一个超时时间 一小时超时
		tm = 1 * time.Hour
	}
	return &notify{
		cc: make(chan int, 1),
		tm: tm,
	}
}

// Signal 发送一个信号
func (s *notify) Signal() {
	s.cc <- 1
}

// 等等一个信息
func (s *notify) Wait() bool {
	timeout := time.After(s.tm)
	select {
	case s.tag = <-s.cc:
	case <-timeout:
		s.tag = -2
	}

	if s.tag == 1 {
		return true
	}

	return false
}

// Stop 发送停止信号
func (s *notify) Stop() {
	s.cc <- -1
}

// Close 关闭
func (s *notify) Close() {
	close(s.cc)
}
