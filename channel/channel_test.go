package channel

import (
	"fmt"
	"testing"
	"time"
)

// 订阅模式
func TestChannel1(t *testing.T) {
	ch := make(chan string, 4)
	go func() {
		str := <-ch
		fmt.Println(str)
	}()
	go func() {
		str := <-ch
		fmt.Println(str)
	}()
	go func() {
		str := <-ch
		fmt.Println(str)
	}()

	ch <- "hello"
	ch <- "hello"
	time.Sleep(time.Second)
}

// 消息队列 方案一
type Broker struct {
	consumers []*Consumer
}

func (b *Broker) Produce(msg string) {
	for _, c := range b.consumers {
		c.ch <- msg
	}
}

func (b *Broker) Subscribe(c *Consumer) {
	b.consumers = append(b.consumers, c)
}

type Consumer struct {
	ch chan string
}

func TestChannel2(t *testing.T) {
	b := &Broker{
		consumers: make([]*Consumer, 0, 10),
	}
	c1 := &Consumer{ch: make(chan string, 1)}
	c2 := &Consumer{ch: make(chan string, 1)}
	b.Subscribe(c1)
	b.Subscribe(c2)
	b.Produce("hello")
	fmt.Println(<-c1.ch)
	fmt.Println(<-c2.ch)
}

// 消息队列 方案二
type Consumer2 struct {
	ch chan string
}

type Broker2 struct {
	ch        chan string
	consumers []func(s string)
}

func (b *Broker2) Produce(msg string) {
	b.ch <- msg
}

func (b *Broker2) Subscribe(consume func(s string)) {
	b.consumers = append(b.consumers, consume)
}

func (b *Broker2) Start() {
	go func() {
		s, ok := <-b.ch
		// 防止接收不到导致goroutine泄漏
		if !ok {
			return
		}
		for _, c := range b.consumers {
			c(s)
		}
	}()
}

// 任务池 方案一
type TaskPool struct {
	ch chan struct{}
}

func NewTaskPool(limit int) *TaskPool {
	t := &TaskPool{
		ch: make(chan struct{}, limit),
	}
	// 提前准备好了令牌
	for i := 0; i < limit; i++ {
		t.ch <- struct{}{}
	}
	return t
}

func (t *TaskPool) Do(f func()) {
	token := <-t.ch
	// 异步执行
	go func() {
		f()
		t.ch <- token
	}()

	// 同步执行
	f()
	t.ch <- token
}

func TestChannel3(t *testing.T) {
	tp := NewTaskPool(2)
	tp.Do(func() {
		time.Sleep(time.Second)
		fmt.Println("task1")
	})
	tp.Do(func() {
		time.Sleep(time.Second)
		fmt.Println("task2")
	})
	tp.Do(func() {
		time.Sleep(time.Second)
		fmt.Println("task3")
	})
}

// 任务池 方案二
type TaskPoolWithCache struct {
	cache chan func()
}

func NewTaskPoolWithCache(limit int, cacheSize int) *TaskPoolWithCache {
	t := &TaskPoolWithCache{
		cache: make(chan func(), cacheSize),
	}
	// 直接把 goroutine开好
	for i := 0; i < limit; i++ {
		go func() {
			for {
				// 在 goroutine 里面不断尝试从 cache 里面拿到任务
				select {
				case task := <-t.cache:
					task()
				}
			}
		}()
	}
	return t
}

func (t *TaskPoolWithCache) Do(f func()) {
	t.cache <- f
}
