package myrunner

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/xxjwxc/public/errors"
)

// 定义新的错误
var (
	ErrRunning = errors.New("任务正在运行")
)

// 任务类型定义
type Task struct {
	Func func(args ...interface{}) // 支持动态参数的任务函数
	Args []interface{}             // 任务执行时的参数
}

// 后台执行任何限时任务，而且我们还可以控制这个执行者，比如强制终止它等
type Runner struct {
	tasks     []Task         //要执行的任务
	complete  chan error     //用于通知任务全部完成
	timeout   time.Duration  //这些任务在多久内完成
	interrupt chan os.Signal //可以控制强制终止的信号
	running   bool           //任务是否正在运行
	mu        sync.Mutex     //互斥锁，确保同一任务不会同时执行
	ticker    *time.Ticker   //定时器，用于定时执行任务
	stopChan  chan struct{}  //用于停止定时任务的通道
}

// 工厂方法
func New(tm time.Duration) *Runner {
	return &Runner{
		complete:  make(chan error), //同步通道，main routine等待，一致要任务完成或者被强制终止
		timeout:   tm,
		interrupt: make(chan os.Signal, 1), //至少接收到一个操作系统的中断信息
		stopChan:  make(chan struct{}),
	}
}

// 添加无参数任务，保持向后兼容
func (r *Runner) Add(tasks ...func()) {
	for _, task := range tasks {
		// 包装为带参数的任务
		r.tasks = append(r.tasks, Task{
			Func: func(args ...interface{}) {
				task()
			},
			Args: nil,
		})
	}
}

// 添加带参数的任务
func (r *Runner) AddWithArgs(f func(args ...interface{}), args ...interface{}) {
	r.tasks = append(r.tasks, Task{
		Func: f,
		Args: args,
	})
}

// 为已添加的任务更新参数
func (r *Runner) UpdateTaskArgs(index int, args ...interface{}) error {
	if index < 0 || index >= len(r.tasks) {
		return errors.New("无效的任务索引")
	}
	r.tasks[index].Args = args
	return nil
}

func (r *Runner) run() error {
	// 执行所有任务
	for _, task := range r.tasks {
		if r.isInterrupt() {
			return ErrInterruput
		}
		task.Func(task.Args...)
	}
	return nil
}

// 检查是否接收到了中断信号
func (r *Runner) isInterrupt() bool {
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}

// 开始执行所有任务，并且监视通道事件
func (r *Runner) Start() error {
	//希望接收哪些系统信号
	signal.Notify(r.interrupt, os.Interrupt) //如果有系统中断的信号，发给r.interrupt

	// 先检查运行状态，如果已经在运行，直接返回ErrRunning
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return ErrRunning
	}
	// 标记为运行状态，防止其他调用进入
	r.running = true
	r.mu.Unlock()

	// 使用局部的complete通道，避免多个调用共享同一个通道
	complete := make(chan error)

	go func() {
		defer func() {
			// 确保在goroutine结束时释放运行状态
			r.mu.Lock()
			r.running = false
			r.mu.Unlock()
		}()
		complete <- r.run()
	}()

	timeoutChan := time.After(r.timeout)
	select {
	case err := <-complete:
		return err
	case <-timeoutChan:
		// 超时后需要释放运行状态
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
		return ErrTimeOut
	}
}

// 立即执行任务
func (r *Runner) RunNow() error {
	//希望接收哪些系统信号
	signal.Notify(r.interrupt, os.Interrupt) //如果有系统中断的信号，发给r.interrupt

	// 先检查运行状态，如果已经在运行，直接返回ErrRunning
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return ErrRunning
	}
	// 标记为运行状态，防止其他调用进入
	r.running = true
	r.mu.Unlock()

	// 使用局部的complete通道，避免多个调用共享同一个通道
	complete := make(chan error)

	go func() {
		defer func() {
			// 确保在goroutine结束时释放运行状态
			r.mu.Lock()
			r.running = false
			r.mu.Unlock()
		}()
		complete <- r.run()
	}()

	timeoutChan := time.After(r.timeout)
	select {
	case err := <-complete:
		return err
	case <-timeoutChan:
		// 超时后需要释放运行状态
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
		return ErrTimeOut
	}
}

// 每天指定时间执行任务
func (r *Runner) StartDaily(hour int) error {
	// 验证小时参数
	if hour < 0 || hour > 23 {
		return errors.New("invalid hour")
	}

	// 计算下一次执行时间
	nextRun := time.Now()
	nextRun = time.Date(nextRun.Year(), nextRun.Month(), nextRun.Day(), hour, 0, 0, 0, nextRun.Location())
	if nextRun.Before(time.Now()) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	// 计算初始延迟
	delay := nextRun.Sub(time.Now())

	// 启动定时任务
	go func() {
		// 初始延迟
		time.Sleep(delay)

		// 执行第一次任务
		r.RunNow()

		// 每天执行一次
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				r.RunNow()
			case <-r.stopChan:
				return
			}
		}
	}()

	return nil
}

// 停止定时任务
func (r *Runner) StopDaily() {
	close(r.stopChan)
}

/*
调用示例:
func main() {
	log.Println("...开始执行任务...")

	timeout := 5 * time.Second
	r := New(timeout)

	// 1. 添加带动态参数的任务a
	taskA := func(args ...interface{}) {
		if len(args) > 0 {
			name, ok := args[0].(string)
			if ok {
				log.Printf("正在执行任务: %s", name)
			}
			if len(args) > 1 {
				count, ok := args[1].(int)
				if ok {
					log.Printf("任务执行次数: %d", count)
				}
			}
		}
		log.Println("任务a执行中...")
		time.Sleep(2 * time.Second)
		log.Println("任务a执行完成")
	}

	// 使用AddWithArgs添加带参数的任务
	r.AddWithArgs(taskA, "测试任务", 1)

	// 2. 定时每天1点执行任务a
	if err := r.StartDaily(1); err != nil {
		log.Println("启动定时任务失败:", err)
		return
	}

	// 3. 手动触发任务a
	go func() {
		// 模拟2秒后手动触发
		time.Sleep(2 * time.Second)
		log.Println("手动触发任务a")
		if err := r.RunNow(); err != nil {
			if err == ErrRunning {
				log.Println("任务a正在执行，跳过本次手动触发")
			} else {
				log.Println("手动执行任务失败:", err)
			}
		}
	}()

	// 4. 更新任务参数并再次触发
	go func() {
		// 模拟5秒后更新参数并触发
		time.Sleep(5 * time.Second)
		log.Println("更新任务a的参数")
		if err := r.UpdateTaskArgs(0, "更新后的任务", 2); err != nil {
			log.Println("更新任务参数失败:", err)
			return
		}
		log.Println("使用新参数手动触发任务a")
		if err := r.RunNow(); err != nil {
			if err == ErrRunning {
				log.Println("任务a正在执行，跳过本次手动触发")
			} else {
				log.Println("手动执行任务失败:", err)
			}
		}
	}()

	// 5. 模拟重复手动触发，测试任务互斥
	go func() {
		// 模拟6秒后再次手动触发，此时任务可能还在执行
		time.Sleep(6 * time.Second)
		log.Println("再次手动触发任务a")
		if err := r.RunNow(); err != nil {
			if err == ErrRunning {
				log.Println("任务a正在执行，跳过本次手动触发")
			} else {
				log.Println("手动执行任务失败:", err)
			}
		}
	}()

	// 运行15秒后退出
	time.Sleep(15 * time.Second)
	log.Println("...程序结束...")
}

*/
