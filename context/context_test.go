package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	withValue := context.WithValue(ctx, "name", "doveOne")

	// Deadline()返回过期时间，如果ok为false，说明没有设置过期时间
	deadline, ok := withValue.Deadline()
	fmt.Println(deadline, ok)

	// 返回一个channel，一般用于监听Context实例的信号
	done := withValue.Done()
	fmt.Println(done)

	fmt.Println("name:", withValue.Value("name"))
	time.Sleep(2 * time.Second)
	err := withValue.Err()
	fmt.Println(err)
}

func MyBusiness() {
	time.Sleep(2 * time.Second)
	fmt.Println("hello")
}

func TestBusinessTimeout(t *testing.T) {
	ctx := context.Background()
	timeout, cancel := context.WithCancel(ctx)
	defer cancel()
	end := make(chan struct{}, 1)
	go func() {
		MyBusiness()
		end <- struct{}{}
	}()
	select {
	case <-timeout.Done():
		fmt.Println("timeout")
	case <-end:
		fmt.Println("business end")
	}
	fmt.Println("over")
}

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer func() {
			fmt.Println("goroutine exit")
		}()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("receive cancel signal!")
				return

			default:
				fmt.Println("default")
				time.Sleep(time.Second)
			}
		}
	}()
	time.Sleep(time.Second)
	defer cancel()
	time.Sleep(2 * time.Second)
}
