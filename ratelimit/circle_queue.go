package ratelimit

import (
	"errors"
	"sync"
	"time"
)

/*
 用环形队列做为底层数据结构来存储用户访问数据,并能实现自动增长以及收缩
*/

// autoGrowCircleQueueInt64 使用切片实现的队列
type autoGrowCircleQueueInt64 struct {
	key interface{}
	// 注意，maxSize比实际存储长度大1
	maxSize int
	// maxSizeTemp与visitorRecord长度相同,visitorRecord长度设计根据实际情况成自动增长
	maxSizeTemp   int
	visitorRecord []int64
	head          int //头
	tail          int //尾
	// 存盘时临时用到的虚拟队列的头和尾
	headForCopy int
	tailForCopy int
	locker      *sync.Mutex
}

// 初始化环形队列,长度超过1023的队列暂时只分配1023的空间
func newAutoGrowCircleQueueInt64(size int) *autoGrowCircleQueueInt64 {
	var c autoGrowCircleQueueInt64
	c.maxSize = size + 1
	if c.maxSize > 1024 {
		c.maxSizeTemp = 1024
	} else {
		c.maxSizeTemp = c.maxSize
	}
	c.visitorRecord = make([]int64, c.maxSizeTemp)
	c.locker = new(sync.Mutex)
	return &c
}

// 队列无人使用时,对于队列实际使用空间长度大于1023的需要对此队列做收缩操作以节省空间
func (q *autoGrowCircleQueueInt64) reSet() {
	q.locker.Lock()
	defer q.locker.Unlock()
	if q.maxSize > 1024 && q.maxSizeTemp > 1024 {
		newVisitorRecord := make([]int64, 1024)
		q.visitorRecord = newVisitorRecord
		q.maxSizeTemp = 1024
		q.head = 0
		q.tail = 0
	}
}

// 队列是否需要扩容
func (q *autoGrowCircleQueueInt64) needGrow() bool {
	if q.maxSizeTemp == q.maxSize {
		return false
	}
	if q.tempQueueIsFull() {
		return true
	}
	return false
}

// 对队列进行扩容操作
func (q *autoGrowCircleQueueInt64) grow() {
	newVisitorRecordLen := len(q.visitorRecord) * 2
	if newVisitorRecordLen > q.maxSize {
		newVisitorRecordLen = q.maxSize
	}
	newVisitorRecord := make([]int64, newVisitorRecordLen)
	// 复制数据
	oldQueueLen := q.tempQueueLen()
	for i := 0; i < oldQueueLen; i++ {
		newVisitorRecord[i] = q.visitorRecord[q.head]
		q.head = (q.head + 1) % q.maxSizeTemp
	}
	// 新旧数据替换
	q.visitorRecord = newVisitorRecord
	q.maxSizeTemp = newVisitorRecordLen
	q.head = 0
	q.tail = oldQueueLen
}

// 访问时间入对列,只用于从本地备份文件加载历史访问数据，本身是线性访问，无并发安全问题
func (q *autoGrowCircleQueueInt64) push(val int64) (err error) {
	q.locker.Lock()
	defer q.locker.Unlock()
	if q.needGrow() {
		q.grow()
	}
	if q.tempQueueIsFull() {
		return errors.New("queue is full")
	}
	q.visitorRecord[q.tail] = val
	q.tail = (q.tail + 1) % q.maxSizeTemp
	return
}

// 访问时间入对列,并发安全,由于不同协程在高并发的时候，极端情况下，也即前后两次访问的时间差，与两协程的系统切换时间非常接近的情况下
// 由调用者自己生成时间容易出现紊乱的情况，所以访问时间只能到这个地方来统一生成，也即有极小的概率，先访问的时间比后访问的时间大
func (q *autoGrowCircleQueueInt64) pushWithConcurrencysafety(defaultExpiration time.Duration) (err error) {
	q.locker.Lock()
	defer q.locker.Unlock()
	if q.needGrow() {
		q.grow()
	}
	if q.tempQueueIsFull() {
		return errors.New("queue is full")
	}
	q.visitorRecord[q.tail] = time.Now().Add(defaultExpiration).UnixNano()
	q.tail = (q.tail + 1) % q.maxSizeTemp
	return
}

// 出对列
func (q *autoGrowCircleQueueInt64) pop() (val int64, err error) {
	q.locker.Lock()
	defer q.locker.Unlock()
	if q.tempQueueIsEmpty() {
		return 0, errors.New("queue is empty")
	}
	val = q.visitorRecord[q.head]
	q.head = (q.head + 1) % q.maxSizeTemp
	return
}

// 用于备份数据的时候，虚拟队列的出队列操作，但实际未进行出队列操作
func (q *autoGrowCircleQueueInt64) tempQueuePopForCopy() (val int64, err error) {
	if q.tempQueueIsEmptyForCopy() {
		return 0, errors.New("queue is empty")
	}
	val = q.visitorRecord[q.headForCopy]
	q.headForCopy = (q.headForCopy + 1) % q.maxSizeTemp
	return
}

// 用于备份数据的时候，判断虚拟队列是否已满
func (q *autoGrowCircleQueueInt64) tempQueueIsFull() bool {
	return (q.tail+1)%q.maxSizeTemp == q.head
}

// 判断队列是否为空
func (q *autoGrowCircleQueueInt64) tempQueueIsEmpty() bool {
	return q.tail == q.head
}

// 用于备份数据的时候，判断虚拟队列是否为空
func (q *autoGrowCircleQueueInt64) tempQueueIsEmptyForCopy() bool {
	return q.tailForCopy == q.headForCopy
}

// 判断队列已使用多少个元素
func (q *autoGrowCircleQueueInt64) usedSize() int {
	return (q.tail + q.maxSizeTemp - q.head) % q.maxSizeTemp
}

// 判断队列中还有多少空间未使用
func (q *autoGrowCircleQueueInt64) tempQueueUnUsedSize() int {
	return q.maxSizeTemp - 1 - q.usedSize()
}

// 判断队列中还有多少空间未使用
func (q *autoGrowCircleQueueInt64) unUsedSize() int {
	return q.maxSize - 1 - ((q.tail + q.maxSizeTemp - q.head) % q.maxSizeTemp)
}

// 队列总的可用空间长度
func (q *autoGrowCircleQueueInt64) tempQueueLen() int {
	return q.maxSizeTemp - 1
}

// 删除过期数据
func (q *autoGrowCircleQueueInt64) deleteExpired(key interface{}) {
	q.locker.Lock()
	defer q.locker.Unlock()
	now := time.Now().UnixNano()
	size := q.usedSize()
	if size == 0 {
		return
	}
	//依次删除过期数据
	for i := 0; i < size; i++ {
		if now > q.visitorRecord[q.head] {
			q.head = (q.head + 1) % q.maxSizeTemp
		} else {
			return
		}
	}
}
