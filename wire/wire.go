//go:build wireinjec
// +build wireinjec

package wire

import (
	"github.com/google/wire"
)

func InitializeBroadCast() BroadCast {
	wire.Build(NewBroadCast, NewChannel, NewMessage)
	return BroadCast{}
}
