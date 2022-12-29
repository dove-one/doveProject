##### wire

```
wire.go 和 wire_gen.go文件头部位置都有一个 +build，不过一个后面是 wireinject，另一个是 !wireinject
+build 其实是 Go 语言的一个特性。类似 C/C++ 的条件编译，在执行 go build 时可传入一些选项，根据这个选项决定某些文件是否编译
wire 工具只会处理有wireinject 的文件，所以我们的 wire.go 文件要加上这个
生成的`wire_gen.go 是给我们来使用的，wire 不需要处理，故有 !wireinject


TODO: wire gen wire.go 时报未识别
xxxx@MacBook-Pro14 wire % wire gen wire.go
wire: /Users/xxxx/doveProject/wire/wire.go:10:28: undeclared name: BroadCast
wire: /Users/xxxx/doveProject/wire/wire.go:11:13: undeclared name: NewBroadCast
wire: /Users/xxxx/doveProject/wire/wire.go:11:27: undeclared name: NewChannel
wire: /Users/xxxx/doveProject/wire/wire.go:11:39: undeclared name: NewMessage
wire: /Users/xxxx/doveProject/wire/wire.go:12:9: undeclared name: BroadCast
wire: generate failed
直接用指令wire成功


1.Provide 构造器


2.injector 注入器
```