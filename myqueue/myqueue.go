package myqueue

import (
	"sync"

	"gopkg.in/eapache/queue.v1"
)

type MyQueue struct {
	lock    sync.Mutex
	popable *sync.Cond
	buffer  *queue.Queue
	closed  bool
}

// 创建
func NewSyncQueue() *MyQueue {
	ch := &MyQueue{
		buffer: queue.New(),
	}
	ch.popable = sync.NewCond(&ch.lock)
	return ch
}

// 取出队列,（阻塞模式）
func (q *MyQueue) Pop() (v interface{}) {
	c := q.popable
	buffer := q.buffer

	q.lock.Lock()
	for buffer.Length() == 0 && !q.closed {
		c.Wait()
	}

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
	}

	q.lock.Unlock()
	return
}

//试着取出队列（非阻塞模式）返回ok == false 表示空
func (q *MyQueue) TryPop() (v interface{}, ok bool) {
	buffer := q.buffer

	q.lock.Lock()

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		ok = true
	} else if q.closed {
		ok = true
	}

	q.lock.Unlock()
	return
}

// 插入队列，非阻塞
func (q *MyQueue) Push(v interface{}) {
	q.lock.Lock()
	if !q.closed {
		q.buffer.Add(v)
		q.popable.Signal()
	}
	q.lock.Unlock()
}

// 获取队列长度
func (q *MyQueue) Len() (l int) {
	q.lock.Lock()
	l = q.buffer.Length()
	q.lock.Unlock()
	return
}

// Close MyQueue
// After close, Pop will return nil without block, and TryPop will return v=nil, ok=True
func (q *MyQueue) Close() {
	q.lock.Lock()
	if !q.closed {
		q.closed = true
		q.popable.Signal()
	}
	q.lock.Unlock()
}
