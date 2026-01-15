package myrunner

import (
	"fmt"
	"testing"
	"time"
)

// 测试Add方法添加无参任务
func TestAdd(t *testing.T) {
	r := New(2 * time.Second)

	// 添加无参任务
	executed := false
	r.Add(func() {
		executed = true
	})

	// 执行任务
	if err := r.RunNow(); err != nil {
		t.Errorf("RunNow failed: %v", err)
	}

	// 验证任务是否执行
	if !executed {
		t.Error("Task was not executed")
	}
}

// 测试AddWithArgs方法添加带参数的任务
func TestAddWithArgs(t *testing.T) {
	r := New(2 * time.Second)

	// 添加带参数的任务
	var result string
	var count int
	r.AddWithArgs(func(args ...interface{}) {
		if len(args) > 0 {
			result, _ = args[0].(string)
		}
		if len(args) > 1 {
			count, _ = args[1].(int)
		}
	}, "test", 123)

	// 执行任务
	if err := r.RunNow(); err != nil {
		t.Errorf("RunNow failed: %v", err)
	}

	// 验证任务是否执行，参数是否正确传递
	if result != "test" {
		t.Errorf("Expected result 'test', got '%s'", result)
	}
	if count != 123 {
		t.Errorf("Expected count 123, got %d", count)
	}
}

// 测试UpdateTaskArgs方法更新任务参数
func TestUpdateTaskArgs(t *testing.T) {
	r := New(2 * time.Second)

	// 添加带参数的任务
	var result string
	r.AddWithArgs(func(args ...interface{}) {
		if len(args) > 0 {
			result, _ = args[0].(string)
		}
	}, "initial")

	// 更新任务参数
	if err := r.UpdateTaskArgs(0, "updated"); err != nil {
		t.Errorf("UpdateTaskArgs failed: %v", err)
	}

	// 执行任务
	if err := r.RunNow(); err != nil {
		t.Errorf("RunNow failed: %v", err)
	}

	// 验证参数是否已更新
	if result != "updated" {
		t.Errorf("Expected result 'updated', got '%s'", result)
	}
}

// 测试任务互斥机制
func TestTaskMutex(t *testing.T) {
	r := New(50 * time.Second)

	// 添加一个长时间运行的任务
	r.Add(func() {
		fmt.Println("任务开始执行")
		time.Sleep(200 * time.Second)
	})

	// 第一次执行任务
	go func() {
		err := r.RunNow()

		fmt.Printf("Expected ErrRunning, got111 %v", err)

	}()

	// 等待任务开始执行
	time.Sleep(100 * time.Millisecond)

	// 第二次执行任务，应该返回ErrRunning
	err := r.RunNow()
	if err != ErrRunning {
		fmt.Printf("Expected ErrRunning, got %v", err)
	}
}

// 测试StartDaily方法
func TestStartDaily(t *testing.T) {
	r := New(2 * time.Second)

	// 验证无效小时参数
	err := r.StartDaily(25)
	if err == nil {
		t.Error("Expected error for invalid hour, got nil")
	}

	// 验证有效小时参数
	err = r.StartDaily(12)
	if err != nil {
		t.Errorf("StartDaily failed for valid hour: %v", err)
	}

	// 停止定时任务
	r.StopDaily()
}

// 测试任务执行超时
func TestTaskTimeout(t *testing.T) {
	r := New(1 * time.Second)

	// 添加一个超时的任务
	r.Add(func() {
		time.Sleep(20 * time.Second)
	})

	// 执行任务，应该返回ErrTimeOut
	err := r.RunNow()
	if err != ErrTimeOut {
		t.Errorf("Expected ErrTimeOut, got %v", err)
	}
}

// 测试UpdateTaskArgs方法的无效索引
func TestUpdateTaskArgsInvalidIndex(t *testing.T) {
	r := New(2 * time.Second)

	// 添加一个任务
	r.Add(func() {})

	// 更新无效索引的任务参数，应该返回错误
	err := r.UpdateTaskArgs(1, "invalid")
	if err == nil {
		t.Error("Expected error for invalid index, got nil")
	}
}
