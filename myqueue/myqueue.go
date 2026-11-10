package myqueue

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/eapache/queue.v1"
)

// MyQueue queue
type MyQueue struct {
	sync.Mutex
	// popable 容量为 1 的唤醒信号通道,语义等价于 sync.Cond.Signal,
	// 但可以配合 select 实现"带超时的等待"。
	// 队列 Close 时会 close 该通道,以此广播唤醒所有等待者。
	popable chan struct{}
	buffer  *queue.Queue
	closed  bool
	count   int32
}

// New 创建
func New() *MyQueue {
	return &MyQueue{
		buffer:  queue.New(),
		popable: make(chan struct{}, 1),
	}
}

// signal 唤醒一个等待者,非阻塞。调用方必须持有锁。
func (q *MyQueue) signal() {
	select {
	case q.popable <- struct{}{}:
	default:
	}
}

// wakeNext 取走数据后若队列仍有剩余,继续唤醒下一个等待者。
// 调用方必须持有锁。
func (q *MyQueue) wakeNext() {
	if q.Len() > 0 {
		q.signal()
	}
}

// Pop 取出队列,（阻塞模式）
func (q *MyQueue) Pop() (v interface{}) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	for q.Len() == 0 && !q.closed {
		q.Mutex.Unlock()
		<-q.popable
		q.Mutex.Lock()
	}

	if q.closed { //已关闭
		return
	}

	if q.Len() > 0 {
		buffer := q.buffer
		v = buffer.Peek()
		buffer.Remove()
		atomic.AddInt32(&q.count, -1)
		q.wakeNext()
	}
	return
}

// TryPop 试着取出队列（非阻塞模式）返回ok == false 表示空
func (q *MyQueue) TryPop() (v interface{}, ok bool) {
	buffer := q.buffer

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if q.Len() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		atomic.AddInt32(&q.count, -1)
		q.wakeNext()
		ok = true
	} else if q.closed {
		ok = true
	}

	return
}

// TryPopTimeout 试着取出队列（阻塞模式+timeout）返回ok == false 表示超时
// 实现上不再"每次调用都新建 goroutine 去等待",避免超时后 goroutine 永久泄漏,
// 也不会出现"持锁向通道发送"导致的死锁。
func (q *MyQueue) TryPopTimeout(tm time.Duration) (v interface{}, ok bool) {
	timer := time.NewTimer(tm)
	defer timer.Stop()

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	for q.Len() == 0 && !q.closed {
		q.Mutex.Unlock()
		select {
		case <-q.popable:
			q.Mutex.Lock()
		case <-timer.C:
			q.Mutex.Lock()
			// 超时:等待期间若已有数据写入或队列已关闭,则继续向下取值
			if q.Len() == 0 && !q.closed {
				return nil, false
			}
		}
	}

	if q.closed { //已关闭
		return nil, true
	}

	if q.Len() > 0 {
		buffer := q.buffer
		v = buffer.Peek()
		buffer.Remove()
		atomic.AddInt32(&q.count, -1)
		q.wakeNext()
	}
	return v, true
}

// Push 插入队列，非阻塞
func (q *MyQueue) Push(v interface{}) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if !q.closed {
		q.buffer.Add(v)
		atomic.AddInt32(&q.count, 1)
		q.signal()
	}
}

// Len 获取队列长度
func (q *MyQueue) Len() int {
	return (int)(atomic.LoadInt32(&q.count))
}

// Close MyQueue
// After close, Pop will return nil without block, and TryPop will return v=nil, ok=True
func (q *MyQueue) Close() {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if !q.closed {
		q.closed = true
		atomic.StoreInt32(&q.count, 0)
		close(q.popable) //广播,唤醒所有阻塞中的等待者
	}
}

// IsClose check is closed
func (q *MyQueue) IsClose() bool {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	return q.closed
}

// Wait 等待队列消费完成
func (q *MyQueue) Wait() {
	for {
		if q.IsClose() || q.Len() == 0 {
			break
		}

		runtime.Gosched() //出让时间片
	}
}
