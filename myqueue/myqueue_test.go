package myqueue

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	que := New()
	for i := 0; i < 10; i++ { //开启20个请求
		que.Push(i)
	}

	go func() {
		for {
			fmt.Println(que.Pop().(int))
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			fmt.Println(que.Pop().(int))
			time.Sleep(1 * time.Second)
		}
	}()

	que.Wait()
	fmt.Println("down")
}

func TestClose(t *testing.T) {
	que := New()
	for i := 0; i < 10; i++ { //开启20个请求
		que.Push(i)
	}

	go func() {
		for {
			v := que.Pop()
			if v != nil {
				fmt.Println(v.(int))
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		for {
			v := que.Pop()
			if v != nil {
				fmt.Println(v.(int))
				time.Sleep(1 * time.Second)
			}
		}
	}()

	que.Close()
	que.Wait()
	fmt.Println("down")
}

func TestTry(t *testing.T) {
	que := New()

	go func() {
		for {
			v, ok := que.TryPop()
			if !ok {
				fmt.Println("no")
				time.Sleep(time.Second / 2)
				runtime.Gosched() //出让时间片
			}

			if v != nil {
				fmt.Println(v.(int))
			}
		}
	}()

	go func() {
		for {
			v, ok := que.TryPop()
			if !ok {
				fmt.Println("no")
				time.Sleep(time.Second / 2)
				runtime.Gosched() //出让时间片
			}

			if v != nil {
				fmt.Println(v.(int))
			}
		}
	}()

	for i := 0; i < 10; i++ { //开启20个请求
		que.Push(i)
		time.Sleep(1 * time.Second)
	}

	que.Wait()
	fmt.Println("down")
}

func TestTimeout(t *testing.T) {
	que := New()
	go func() {
		for i := 0; i < 10; i++ { //开启20个请求
			time.Sleep(1 * time.Second)
			que.Push(i)

		}
	}()

	go func() {
		for {
			b, ok := que.TryPopTimeout(100 * time.Microsecond)
			if ok {
				fmt.Println(b.(int))
			} else {
				fmt.Println("time out")
			}
		}
	}()

	time.Sleep(200 * time.Second)
	que.Wait()
	fmt.Println("down")
}
