package reflect

import (
	"reflect"
	"testing"
)

type DoveOne struct {
	Name    string
	Age     int32
	Alias   []byte
	Address string
}

func BenchmarkFieldIndexOrName(b *testing.B) {
	tm := DoveOne{}
	val := reflect.ValueOf(tm)
	b.Run("by index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = val.Field(0)
		}
	})

	b.Run("by name", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = val.FieldByName("Age")
		}
	})

}
