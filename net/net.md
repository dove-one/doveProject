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
这些方法都是返回Listener的具体类，如TCPListener。一般用Listen就可以
除非你要依赖于具体的网络协议特性
网络通信用TCP还是UDP是一个影响巨大的事情，一般确认了就不会改。

```