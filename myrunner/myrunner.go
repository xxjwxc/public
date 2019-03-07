package myrunner

import (
	"os"
	"os/signal"
	"time"
)

//后台执行任何限时任务，而且我们还可以控制这个执行者，比如强制终止它等
type Runner struct {
	tasks     []func()         //要执行的任务
	complete  chan error       //用于通知任务全部完成
	timeout   <-chan time.Time //这些任务在多久内完成 只能接收
	interrupt chan os.Signal   //可以控制强制终止的信号
}

//工厂方法
func New(tm time.Duration) *Runner {
	return &Runner{
		complete:  make(chan error), //同步通道，main routine等待，一致要任务完成或者被强制终止
		timeout:   time.After(tm),
		interrupt: make(chan os.Signal, 1), //至少接收到一个操作系统的中断信息
	}
}

//
func (r *Runner) Add(tasks ...func()) {
	r.tasks = append(r.tasks, tasks...)
}

//
func (r *Runner) run() error {
	for _, task := range r.tasks {
		if r.isInterrupt() {
			return ErrInterruput
		}
		task()
	}
	return nil
}

//检查是否接收到了中断信号
func (r *Runner) isInterrupt() bool {
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}

//开始执行所有任务，并且监视通道事件
func (r *Runner) Start() error {
	//希望接收哪些系统信号
	signal.Notify(r.interrupt, os.Interrupt) //如果有系统中断的信号，发给r.interrupt

	go func() {
		r.complete <- r.run()
	}()

	select {
	case err := <-r.complete:
		return err
	case <-r.timeout:
		return ErrTimeOut
	}
}

/*
调用示例:
func main() {
	log.Println("...开始执行任务...")

	timeout := 2 * time.Second
	r := New(timeout)

	r.Add(createTask(0), createTask(1), createTask(2))

	if err := r.Start(); err != nil {
		switch err {
		case ErrTimeOut:
			log.Println(err)
			os.Exit(1) //退出
		case ErrInterruput:
			log.Println(err)
			os.Exit(2)
		default:
			break
		}
	}
	log.Println("...任务执行结束...")
}

func createTask(param int) func() {
	return func() {
		log.Printf("正在执行任务%d", param)
		time.Sleep(time.Duration(param) * time.Second)
	}
}

*/
