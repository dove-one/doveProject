package unsafe

import (
	"fmt"
	"reflect"
	"testing"
)

type DoveOne struct {
	Name    string
	Age     int32
	Alias   []byte
	Address string
}

func TestUnsafe(t *testing.T) {
	PrintFieldOffset(DoveOne{})
}

// PrintFieldOffset 用来打印字段偏移量
// 用于研究内存布局
// 只接受结构体作为输入
func PrintFieldOffset(entity any) {
	typ := reflect.TypeOf(entity)
	for i := 0; i < typ.NumField(); i++ {
		fd := typ.Field(i)
		fmt.Println(fd.Name, fd.Offset)
	}
}
