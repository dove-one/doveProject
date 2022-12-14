package v1

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

// channelPool 实现了基于缓冲通道的 Pool 接口
type channelPool struct {
	// 用于储存我们的 net.Conn 连接
	mu    sync.RWMutex
	conns chan net.Conn

	// net.Conn 生成器
	factory Factory
}

// Factory 是创建新连接的方法
type Factory func() (net.Conn, error)

// NewChannelPool 返回一个基于缓冲通道的新连接池，有一个初始容量和最大容量
func NewChannelPool(initialCap, maxCap int, factory Factory) (Pool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	c := &channelPool{
		conns:   make(chan net.Conn, maxCap),
		factory: factory,
	}

	for i := 0; i < initialCap; i++ {
		conn, err := factory()
		if err != nil {
			// 如果创建初始连接失败，就关闭连接池
			c.Close()
			return nil, fmt.Errorf("factory is not able to fill pool: %s", err)
		}
		c.conns <- conn
	}
	return c, nil
}

func (c *channelPool) getConnsAndFactory() (chan net.Conn, Factory) {
	c.mu.RLock()
	defer c.mu.Unlock()
	factory := c.factory
	conns := c.conns
	return conns, factory

}

func (c *channelPool) Get() (net.Conn, error) {
	conns, factory := c.getConnsAndFactory()
	if conns == nil {
		return nil, ErrClosed
	}
	select {
	case conn := <-conns:
		if conns == nil {
			return nil, ErrClosed
		}
		return c.wrapConn(conn), nil
	default:
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		return c.wrapConn(conn), nil
	}
}

// put() 把连接放回连接池中
func (c *channelPool) put(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	c.mu.RLock()
	defer c.mu.Unlock()

	if c.conns == nil {
		// 池已经关闭，关闭已通过的连接
		return conn.Close()
	}

	// 将资源放回池中。如果池子满了，将会阻塞并执行默认的情况
	select {
	case c.conns <- conn:
		return nil
	default:
		// 池满了，关闭已通过的连接
		return conn.Close()
	}

}

func (c *channelPool) Close() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
}

func (c *channelPool) Len() int {
	conns, _ := c.getConnsAndFactory()
	return len(conns)
}
