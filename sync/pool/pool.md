##### Pool

```
保存和复用临时对象，减少内存分配，降低 GC 压力。
并发安全

一般情况下，如果要考虑缓存资源，比如创建好的对象，那么可以使用sync.pool
1.sync.Pool会先查看自己是否有资源，有则直接返回
2.没有则创建一个新的
3.sync.pool会在GC的时候释放缓存的资源

一般用sync.Pool都是为了复用内存：
1.它减少了内存分配，也减轻了GC压力(最主要)
2.少消耗CPU资源(内存分配和GC都是CPU密集操作)

TODO：看TLB(thread-local=buffer)方案，go使用的go-local方案(窃取)

*Go的设计：
·每个P一个poolLocal对象
·每个poolLocal有一个private和shared
·shared指向的是一个poolChain。poolChain的数据会被别的P给偷走
·poolChain是一个链表+ring buffer的双重结构
 · 从整体上来说，他是一个双向链表
 · 从单个节点来说，它指向一个ring buffer。后一个节点的ring buffer都是前一个节点的两倍
 
ring buffer优势(实际上也可以说是数组的优势):
1.一次性分配后内存，循环利用
2.对缓存友好

pool.Get()步骤：
1.看private可不可用，可用就直接返回
2.不可用则从自己的poolChain里面尝试获取一个
  · 从头开始找。注意，头指向的是最近创建的ring buffer
  · 从队头往队尾找
3.找不到则尝试从别的P里面投以个出来。偷的过程就是全局并发，因为理论上，其他P都可能恰好一起来偷了
  · 偷是从队尾偷的
4.如果偷也偷不到，那么就会去找缓刑(victim)的
5.连缓刑的也没有，那就去创建一个新的
注：45先偷后victim顺序是因为sync.Pool希望victime里面的对象尽可能被回收掉。

pool.Put()步骤
1.private要是没放东西，就直接放private
2.否则，准备放poolChain
 ·如果poolChain的HEAD还没创建，就创建一个HEAD，然后创建一个8容量的ring buffer，把数据丢过去
 ·如果poolChain的HEAD指向的ring buffer没满，则丢过去ring buffer
 ·如果poolChain的HEAD指向的ring buffer已经满了，就创建一个新的节点，并且创建一个两倍容量的ring buffer，把数据丢过去
 
Go的sync.Pool纯粹一览于GC，用户完全没办法手动控制。
sync.Pool的核心机制是依赖于两个：
 ·locals
 ·victim：缓刑
 
GC过程：
1.locals会被挪过去变成victim
2.victim会被直接回收掉

复活：如果victim的对象再次被使用，那么它就会被丢回去locals，逃过了下一轮被GC回收掉的命运
优点：防止GC引起性能抖动

poolLocal和false sharding
每个poolLocal都有一个pad字段，是用于将poolLocal所占用的内存补齐到128的整数倍。
在并发话题下：所有的对齐基本上都是为了独占CPU高速缓存的CacheLine

*注：封装pool要注意GC

https://github.com/valyala/bytebufferpool 字节缓存池

sync.Pool面试要点
1.sync.Pool和GC的关系：数据默认在local里面，GC的时候会被挪过去victim里面。如果这时候有P用了victim的数据，那么数据会被放回去local里面。
2.poolChain的设计：核心在于理解poolChain是一个双向链表加ring buffer的双重结构

question
1.什么时候P会用victim的数据：偷都偷不到的时候。
2.为什么Go会设计这种结构？一个全局共享队列不好吗？这个问题结合TLB来回答，TLB解决全局锁竞争的方案，Go结合自身P这么一个优势，设计出来的。
3.窃取：结合GMP调度里面的工作窃取，原理一样
4.使用sync.Pool有什么注意点(缺点，优点)？高版本的Go里的sync.Pool没特别大的缺点，硬要说就是内存使用量不可控，以及GC之后即便可以用victim，Get的速率还是要差点。
```

##### 