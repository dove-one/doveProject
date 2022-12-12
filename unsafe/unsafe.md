##### unsafe

```
内存对齐 字段偏移
ptr是字段偏移量：
ptr = 结构体起始地址 + 字段偏移量
见之前的笔记 待补

unsafe.Pointer和uintptr
· unsafe.Pointer：是Go层面的指针，GC会维护unsafe.Pointer的值
· uintptr:直接就是一个数字，代表的是一个内存地址

unsafe.Pointer和GC
假如说GC前一个unsafe.Pointer代表对象的指针，它此时指向的地址是0xAAAA
如果发送了GC，GC之后这个对象依旧存活，但是此时这个对象被复制过去了另外一个位置
(Go GC算法是标记-复制)。那么此时代表对象的unsafe.Pointer会被GC修正，指向新的地址0xAABB

uintptr使用误区
如果使用uintptr来保存对象的起始地址，那么如果发送GC了，原本的代码会直接崩溃。
例如在GC前，计算到的entityAddr = 0xAAAA,那么GC后因为复制的原因，实际上的地址变成了0xAABB
因为GC不会维护uintptr变量，所以entityAddr还是0xAAAA，这个时候再用0xAAAA作为起始地址去访问字段，

但是uintptr可以用于表达相对的量。
例如字段偏移量。这个字段的偏移量是不管怎么GC都不会变的。
如果怕出错，那么就只在进行地址运算的时候使用uintptr，其他时候都用unsafe.Pointer

unsafe面试
· uintptr和unsafe.Pointer的区别：前者代表的是一个具体的地址，后者代表的是一个逻辑上的指针。
 后者在GC等情况下，go runtime会帮你调整，使其永远指向真是存放对象的地址，
· Go对象是怎么对齐的？按照字长。 可能要你手动演示如何对齐和计算对象的大小(之后看一遍笔记)
· 怎么计算对象地址？对象的起始地址是通过反射来获取，对象内部字段的地址是通过起始地址+字段偏移量来计算。
· unsafe为什么比反射高效？可以简单认为反射帮我们封装了很多unsafe的操作，所以我们直接使用unsafe绕过
 了这些封装的开销。

```
