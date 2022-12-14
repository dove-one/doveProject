package v1

import (
	"errors"
	"net"
)

var (
	// ErrClosed 是连接池关闭 pool.Close() 的时候发生的错误
	ErrClosed = errors.New("pool is closed")
)

// Pool 接口描述了一个连接池的实现。
// 一个连接池应该有最大容量。
// 理想的连接池是线程安全的且易于使用
type Pool interface {
	// Get 返回一个新的连接
	// 关闭连接后会将它放回池中
	// 当连接池被破坏或者满的时候关闭会被当作一个错误
	Get() (net.Conn, error)

	// Close 关闭了连接池和它所有的连接
	Close()

	// Len 返回连接池的当前连接数
	Len() int
}
