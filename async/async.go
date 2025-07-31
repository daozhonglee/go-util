// Package async 提供异步操作和安全协程工具
package async

import (
	"fmt"
	"runtime/debug"
	"time"
)

type Future struct {
	ch chan interface{}
}

func (f *Future) Get() interface{} {
	return <-f.ch
}

// Go 异步执行函数，返回Future对象
func Go(fun func() interface{}) *Future {
	future := &Future{ch: make(chan interface{}, 1)}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				close(future.ch)
				s := string(debug.Stack())
				fmt.Printf("[Async.Go] func panic, err = %s, stack = %s", err, s)
			}
		}()

		future.ch <- fun()
	}()

	return future
}

// Timeout 带超时的函数执行，超时返回nil
func Timeout(fun func() interface{}, timeout time.Duration) interface{} {
	// 添加缓冲区，防止超时返回后，函数结果写入channel阻塞，内存无法回收
	ch := make(chan interface{}, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				close(ch)
				s := string(debug.Stack())
				fmt.Printf("[Async.Timeout] func panic, err = %s, stack = %s", err, s)
			}
		}()

		ch <- fun()
	}()

	select {
	case <-time.After(timeout):
		return nil
	case data := <-ch:
		return data
	}
}

// GoTimeout 异步执行带超时的函数
func GoTimeout(fun func() interface{}, timeout time.Duration) *Future {
	nf := func() interface{} {
		return Timeout(fun, timeout)
	}
	return Go(nf)
}

// Sync 同步执行函数，立即返回结果
func Sync(fun func() interface{}) *Future {
	future := &Future{ch: make(chan interface{}, 1)} //同步channel需要缓冲区
	future.ch <- fun()
	return future
}

// Safe 安全启动goroutine，自动捕获panic
func Safe(fun func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				s := string(debug.Stack())
				fmt.Printf("[Async.Safe] func panic, err = %s, stack = %s", err, s)
			}
		}()

		fun()
	}()
}
