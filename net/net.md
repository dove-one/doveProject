##### net

```
在net里面，最重要的两个调用：
· Listen(network, addr sting):监听某个端口，等待客户端连接
· Dial(network, addr string):拨号，其实也就是连上某个服务端

基本分为两大阶段。
创建连接阶段：
1.服务端开始监听一个端口
2.客户端拨通服务器端，两者协商创建连接(TCP)

通信阶段：
1.客户端不断发送请求
2.服务端读取请求
3.服务端处理请求
4.服务端写回响应

net.Listen
Listen是监听一个端口，准备读取数据。
· ListenTcp
· ListenUDP
· ListenIP
· ListenUnix
这些方法都是返回Listener的具体类，如TCPListener。一般用Listen就可以
除非你要依赖于具体的网络协议特性
网络通信用TCP还是UDP是一个影响巨大的事情，一般确认了就不会改。

net.Dial
Dial是指创建一个连接，连上远端的服务器。
· DialIP
· DialTCP
· DialUDP
· DialUnix
· DialTimeout
只有DialTimeout稍微特殊一点，它多了一个超时参数。
类似于Listen，建议直接使用DialTimeout，因为设置超时可以避免一直阻塞

面试要点
1.网络的基础知识，包含TCP和UDP的基础知识
 三次握手和四次挥手
2.如何利用Go写一个简单的TCP服务(之后补demo)
3.记住goroutine和连接的关系，可以在不同的环节使用不同的goroutine，以充分利用
 TCP的全双工通信
```

连接池
```
连接池
· 要发起系统调用
· TCP要完成三次握手
· 高并发的情况下，可能耗尽文件描述符
连接池就是为了复用这些创建好的连接

通常，连接池有几个参数：最小连接数，空闲连接数，最大连接数
v1版本见 doveProject/connect_pool
```