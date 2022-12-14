package v2

import (
	"errors"
)

var (
	// ErrClosed 是连接池关闭 pool.Close() 的时候发生的错误
	ErrClosed = errors.New("pool is closed")
)

// Pool 接口描述了一个连接池的实现。
// 一个连接池应该有最大容量。
// 理想的连接池是线程安全的且易于使用
type Pool interface {
	Get() (any, error)

	Put(any) error

	Close(any) error

	Len() int

	Release()
}
