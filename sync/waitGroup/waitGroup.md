##### WaitGroup

```
WaitGroup是用于同步多个goroutine之间的。

常见场景是我们会把任务拆分给多个goroutine并行完成。
在完成之后需要合并这些任务的结果，或者需要等腰所有小任务都完成之后才能进入下一步。

WaitGroup是用于同步多个goroutine之间工作的。
· 要在开启goroutine之前先加1
· 每一个小任务完成就减1
· 调用Wait方法来等待所有子任务完成

容易犯错的地方是+1和-1不匹配(非常不好测试)：
· 加多了导致Wait一直阻塞，引起goroutine泄漏
· 减多了直接就panic

type WaitGroup struct {
	noCopy noCopy
	
	state1 uint64
	// 64位值：高32位是计数器，低32位是等待器计数。协程xx
	// counter + waiter (=statep)
	state2 uint32
	// sema锁 装载等待的协程
}

WaitGroup
· state1:在64位下，高32位记录了还有多少任务在运行；
         低32位记录了有多少goroutine在等Wait()方法返回
· state2:信号量，用于挂起或者唤醒gorountine,约等于Mutex里面的sema字段
         (要注意很想对比)
         
本质上，WaitGroup是一个无锁实现，严重依赖于CAS对state1的操作

WaitGroup细节
· Add：看上去就是state1的高32位自增1，原子操作
· Done：看上去就是state1的高32位自减1，原子操作，然后看看是不是要唤醒等待goroutine
· Wait：看上去就是state1的低32位自增1，同时利用state2和runtime_Semacquire调用把当前gorountine挂起
其实Done就是相当于Add(-1)

与errgroup对比
可以认为errgroup.Group是对WaitGroup的封装
· 首先需要引入golang.org/x/sync依赖
· errgroup.Group会帮我们保持进行中任务计数
· 任何一个任务返回error，Wait方法就会返回error
大多数情况下，随便选择哪个都可以，差异不大。
```

##### 