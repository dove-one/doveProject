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
```

##### 