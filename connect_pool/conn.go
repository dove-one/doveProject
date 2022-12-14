package connect_pool

import (
	"net"
	"sync"
)

// PoolConn 是对 net.Conn 的封装，用于修改 net.Conn Close()
type PoolConn struct {
	net.Conn
	mu       sync.RWMutex
	c        *channelPool
	unusable bool
}

// Close() 把连接放回连接池中，而不是关闭它
func (p *PoolConn) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
		return nil
	}
	return p.c.put(p.Conn)
}

// MarkUnusable() 标记连接不再使用，让连接池关闭它，而不是把它放回连接池
func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.unusable = true
}

// wrapConn()把一个标准的 net.Conn 封装成 poolConn net.conn
func (c *channelPool) wrapConn(conn net.Conn) net.Conn {
	p := &PoolConn{
		c: c,
	}
	p.Conn = conn
	return p
}
