#### context

```
context线程安全 why？

context包核心api
·context.WithValue: 设置键值对，并且返回一个新的context实例
·context.WithCancel
·context.WithDeadline
·context.WithTimeout:三者都返回一个可取消的context实例，和取消函数(defer cancel())
(1安全传递数据，234用于控制链路)

*注意：context实例是不可变的，每一次都是新创建的

Context接口核心API
·Deadline：返回过期时间，如果ok为false，说明没有设置过期时间。不常用
·Done：返回一个channel，一般用于监听Context实例的信号，比如说过期，或者正常关闭。常用
·Err：返回一个错误用于表达Context发生了什么。Canceled => 正常关闭，DeadlineExceeded => 过期超时。比较常用
·context.Value：取值。非常常用

valueCtx 用于存储 key-value 数据，特点：
 1.典型的装饰器模式：在已有Context的基础上附加一个存储key-value的功能 
 2.只能存储一个key, val
  - 为什么不用 map?
map要求key是comparable的，而我们可能用不是comparable的key。context包的设计理念就是将Context设计成不可变

*TODO 看设计模式装饰器模式和适配器模式，debug断点看value方法的内部逻辑

*case <-ctx.Done() 多用于超时控制
context.Context普遍会和select-case一起使用

cancelCtx:父亲取消的时候，会把儿子取消掉
timerCtx:
·WithTimeout和WithDeadline本质一样
·WithDeadline里面，在创建timerCtx的时候利用time.AfterFunc来实现超时)

注：难点propagateCancel(parent, c)将儿子捆绑到父亲context里面，这里不必细看知道方法作用

使用注意事项
1.一般只用作方法参数，而且是作为第一个参数
2.所有公有方法，除非是util,helper之类的方法，否则都加上context参数
3.不要用做结构体字段，除非你的结构体本身也是表达一个上下文的概念

context面试要点
·context.Context使用原理：上下文传递和超时控制
·context.Context原理：
 ·父亲如何控制儿子：通过儿子主动加入到父亲的children里面，父亲只需要遍历就可以
 ·valueCtx和timeCtx的原理
 
 *TODO：选项模式 options模式整理文章
 
 装饰器和适配器模式整理文章上github 看优雅退出的ppt 再看文章
```

