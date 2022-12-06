package once

import (
	"fmt"
	"sync"
	"testing"
)

type OnceClose struct {
	close sync.Once
}

func (o *OnceClose) Close1() error {
	o.close.Do(func() {
		fmt.Println("close")
	})
	return nil
}

// *必须指针，不加指针每次调用 结构体会被复制一次
func (o OnceClose) Close2() error {
	o.close.Do(func() {
		fmt.Println("close")
	})
	return nil
}

func TestOnce(t *testing.T) {
	o := &OnceClose{}
	//o.Close1()
	//o.Close1()
	o.Close2()
	o.Close2()
}

func init() {
	// 在这里的动作，肯定执行一次
	// 常做包级别的初始化操作
}
