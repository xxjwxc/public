package myqueue

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPopPush(t *testing.T) {
	q := New()
	defer q.Close()

	go func() {
		time.Sleep(20 * time.Millisecond)
		q.Push(1)
	}()

	done := make(chan interface{}, 1)
	go func() {
		done <- q.Pop()
	}()

	select {
	case v := <-done:
		if v != 1 {
			t.Fatalf("期望取到 1,实际 %v", v)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Pop 未被 Push 唤醒")
	}
}

func TestTryPop(t *testing.T) {
	q := New()

	if v, ok := q.TryPop(); ok || v != nil {
		t.Fatalf("空队列 TryPop 应返回 (nil,false),实际 (%v,%v)", v, ok)
	}

	q.Push("a")
	if v, ok := q.TryPop(); !ok || v != "a" {
		t.Fatalf("期望取到 a,实际 (%v,%v)", v, ok)
	}

	q.Close()
	if v, ok := q.TryPop(); !ok || v != nil {
		t.Fatalf("关闭后 TryPop 应返回 (nil,true),实际 (%v,%v)", v, ok)
	}
}

func TestTryPopTimeout(t *testing.T) {
	q := New()
	defer q.Close()

	// 1. 空队列应超时,且不提前返回
	start := time.Now()
	v, ok := q.TryPopTimeout(50 * time.Millisecond)
	if ok || v != nil {
		t.Fatalf("空队列应超时,实际 (%v,%v)", v, ok)
	}
	if d := time.Since(start); d < 40*time.Millisecond {
		t.Fatalf("超时时间明显偏短: %v", d)
	}

	// 2. 已有数据应立即返回
	q.Push(42)
	if v, ok = q.TryPopTimeout(time.Second); !ok || v != 42 {
		t.Fatalf("期望取到 42,实际 (%v,%v)", v, ok)
	}

	// 3. 等待期间被 Push 唤醒,应早于超时返回
	go func() {
		time.Sleep(20 * time.Millisecond)
		q.Push(7)
	}()
	if v, ok = q.TryPopTimeout(2 * time.Second); !ok || v != 7 {
		t.Fatalf("期望在超时前取到 7,实际 (%v,%v)", v, ok)
	}

	// 4. 反复超时不应泄漏 goroutine
	base := runtime.NumGoroutine()
	for i := 0; i < 50; i++ {
		q.TryPopTimeout(time.Millisecond)
	}
	// 给已退出的 goroutine 一点调度时间
	time.Sleep(50 * time.Millisecond)
	if n := runtime.NumGoroutine(); n > base+5 {
		t.Fatalf("TryPopTimeout 疑似泄漏 goroutine: 之前 %d,现在 %d", base, n)
	}
}

func TestCloseWakeupAll(t *testing.T) {
	q := New()

	const waiters = 3
	var wg sync.WaitGroup
	wg.Add(waiters)
	for i := 0; i < waiters; i++ {
		go func() {
			defer wg.Done()
			if v := q.Pop(); v != nil {
				t.Errorf("关闭后 Pop 应返回 nil,实际 %v", v)
			}
		}()
	}

	time.Sleep(50 * time.Millisecond) // 让所有 goroutine 进入等待
	q.Close()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Close 未能唤醒所有阻塞的 Pop")
	}

	if !q.IsClose() {
		t.Error("IsClose 应返回 true")
	}
}

func TestMultipleWaitersGetEach(t *testing.T) {
	q := New()
	defer q.Close()

	const n = 3
	got := make(chan int, n)
	for i := 0; i < n; i++ {
		go func() {
			if v := q.Pop(); v != nil {
				got <- v.(int)
			}
		}()
	}
	time.Sleep(50 * time.Millisecond) // 让所有 goroutine 进入等待

	for i := 0; i < n; i++ {
		q.Push(i)
	}

	seen := make(map[int]bool)
	for i := 0; i < n; i++ {
		select {
		case v := <-got:
			seen[v] = true
		case <-time.After(2 * time.Second):
			t.Fatalf("只收到 %d 个结果就超时: %v", len(seen), seen)
		}
	}
	if len(seen) != n {
		t.Fatalf("期望收到 %d 个不同值,实际 %d", n, len(seen))
	}
}

func TestConcurrentPushPop(t *testing.T) {
	q := New()
	const total = 2000
	const workers = 8

	var consumed int32
	var sum int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				v := q.Pop()
				if v == nil {
					return
				}
				atomic.AddInt32(&consumed, 1)
				atomic.AddInt64(&sum, int64(v.(int)))
			}
		}()
	}

	for i := 0; i < total; i++ {
		q.Push(i)
	}

	// 等全部消费完再 Close,否则剩余数据会被丢弃
	deadline := time.Now().Add(3 * time.Second)
	for atomic.LoadInt32(&consumed) < total && time.Now().Before(deadline) {
		runtime.Gosched()
	}
	q.Close()
	wg.Wait()

	if got := atomic.LoadInt32(&consumed); got != total {
		t.Fatalf("期望消费 %d 个,实际 %d", total, got)
	}
	want := int64(total * (total - 1) / 2)
	if got := atomic.LoadInt64(&sum); got != want {
		t.Fatalf("求和期望 %d,实际 %d", want, got)
	}
}

func TestWait(t *testing.T) {
	q := New()
	for i := 0; i < 5; i++ {
		q.Push(i)
	}

	go func() {
		for i := 0; i < 5; i++ {
			q.Pop()
		}
	}()

	done := make(chan struct{})
	go func() {
		q.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait 未在队列消费完后返回")
	}
}
