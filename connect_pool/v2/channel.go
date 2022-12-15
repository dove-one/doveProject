package v2

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	//ErrMaxActiveConnReached 连接池超限
	ErrMaxActiveConnReached = errors.New("MaxActiveConnReached")
)

// Config 连接池相关配置
type Config struct {
	InitialCap  int
	MaxCap      int
	MaxIdle     int
	Factory     func() (any, error)
	Close       func(any) error
	Ping        func(any) error
	IdleTimeout time.Duration
}

type connReq struct {
	idleConn *idleConn
}

// channelPool 实现了基于缓冲通道的 Pool 接口
type channelPool struct {
	mu           sync.RWMutex
	conns        chan *idleConn
	factory      func() (any, error)
	close        func(any) error
	ping         func(any) error
	idleTimeout  time.Duration
	waitTimeout  time.Duration
	maxActive    int
	openingConns int
	connReqs     []chan connReq
}

type idleConn struct {
	conn any
	t    time.Time
}

// NewChannelPool 返回一个基于缓冲通道的新连接池
func NewChannelPool(p *Config) (Pool, error) {
	if !(p.InitialCap <= p.MaxIdle && p.MaxCap >= p.MaxIdle && p.InitialCap >= 0) {
		return nil, errors.New("invalid capacity settings")
	}
	if p.Factory == nil {
		return nil, errors.New("invalid factory func settings")
	}
	if p.Close == nil {
		return nil, errors.New("invalid close func settings")
	}
	c := &channelPool{
		mu:           sync.RWMutex{},
		conns:        make(chan *idleConn, p.MaxIdle),
		factory:      p.Factory,
		close:        p.Close,
		idleTimeout:  p.IdleTimeout,
		maxActive:    p.MaxCap,
		openingConns: p.InitialCap,
	}
	if p.Ping != nil {
		c.ping = p.Ping
	}
	for i := 0; i < p.InitialCap; i++ {
		conn, err := c.factory()
		if err != nil {
			// 如果创建初始连接失败，就关闭连接池
			c.Release()
			return nil, fmt.Errorf("factory is not able to fill pool: %s", err)
		}
		c.conns <- &idleConn{
			conn: conn,
			t:    time.Now(),
		}
	}
	return c, nil
}

// getConns 获取所有连接
func (c *channelPool) getConns() chan *idleConn {
	c.mu.Lock()
	defer c.mu.Unlock()
	conns := c.conns
	return conns
}

// Get 从pool中获取一个连接
func (c *channelPool) Get() (any, error) {
	conns := c.conns
	if conns == nil {
		return nil, ErrClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, ErrClosed
			}
			// 判断是否超时，如果超时则丢弃
			if c.idleTimeout > 0 {
				if wrapConn.t.Add(c.idleTimeout).Before(time.Now()) {
					// 丢弃并关闭该连接
					c.Close(wrapConn.conn)
					continue
				}
			}
			// 判断是否失效，如果失效则丢弃；若用户没有设定 ping 方法，则不检查
			if c.ping != nil {
				if c.Ping(wrapConn.conn) != nil {
					c.Close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			c.mu.Lock()
			if c.openingConns >= c.maxActive {
				req := make(chan connReq, 1)
				c.connReqs = append(c.connReqs, req)
				c.mu.Unlock()
				ret, ok := <-req
				if !ok {
					return nil, ErrMaxActiveConnReached
				}
				if c.idleTimeout > 0 {
					if ret.idleConn.t.Add(c.idleTimeout).Before(time.Now()) {
						// 丢弃并关闭该连接
						c.Close(ret.idleConn.conn)
						continue
					}
				}
				return ret.idleConn.conn, nil
			}
			if c.factory == nil {
				c.mu.Unlock()
				return nil, ErrClosed
			}
			conn, err := c.factory()
			if err != nil {
				c.mu.Unlock()
				return nil, err
			}
			c.openingConns++
			c.mu.Unlock()
			return conn, nil
		}
	}
}

// Put 将连接放回池中
func (c *channelPool) Put(conn any) error {
	if conn == nil {
		return errors.New("connection is nil.rejecting")
	}

	c.mu.Lock()

	if c.conns == nil {
		c.mu.Unlock()
		return c.Close(conn)
	}

	if len(c.connReqs) > 0 {
		l := len(c.connReqs)
		req := c.connReqs[0]
		copy(c.connReqs, c.connReqs[1:])
		c.connReqs = c.connReqs[:l-1]
		req <- connReq{
			idleConn: &idleConn{
				conn: conn,
				t:    time.Now(),
			},
		}
		c.mu.Unlock()
		return nil
	} else {
		select {
		case c.conns <- &idleConn{conn: conn, t: time.Now()}:
			c.mu.Unlock()
			return nil
		default:
			c.mu.Unlock()
			// 连接池已满，直接关闭该连接
			return c.Close(conn)
		}
	}
}

// Close 关闭单条连接
func (c *channelPool) Close(conn any) error {
	if conn == nil {
		return errors.New("connection is nil.rejecting")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.close == nil {
		return nil
	}
	c.openingConns--
	return c.close(conn)
}

// Ping 检查单条连接是否有效
func (c *channelPool) Ping(conn any) error {
	if conn == nil {
		return errors.New("connection is nil.rejecting")
	}
	return c.ping(conn)
}

// Len 连接池中已有的连接
func (c *channelPool) Len() int {
	return len(c.getConns())
}

// Release 释放连接池中所有连接
func (c *channelPool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.ping = nil
	closeFunc := c.close
	c.close = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}
	close(conns)
	for wrapConn := range conns {
		closeFunc(wrapConn)
	}
}
