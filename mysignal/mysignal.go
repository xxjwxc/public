package mysignal

import (
	"os"
	"os/signal"
	"syscall"
)

type notify struct {
	cc chan os.Signal
}

// New new signal
func New() *notify {
	return &notify{
		cc: make(chan os.Signal, 1),
	}
}

func (s *notify) Wait() {
	signal.Notify(s.cc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	// wait on kill signal
	case <-s.cc:
	}
}

// NotifyStop 发送停止信号
func (s *notify) NotifyStop() {
	s.cc <- syscall.SIGINT
}
