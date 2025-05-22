package enums

import (
	"fmt"

	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/types/enumx"
)

// Print 打印枚举值
func Print(color fmtx.Color, enum *enumx.Enum[string, string]) {
	for _, k := range enum.Keys() {
		fmt.Printf("%-30s %s\n", color.String(k), enum.Get(k))
	}
}

// PrintDesc 反向打印枚举值
func PrintDesc(color fmtx.Color, enum *enumx.Enum[string, string]) {
	l := len(enum.Keys())
	for i := l; i > 0; i-- {
		key := enum.Keys()[i-1]
		value := enum.Get(key)
		fmt.Printf("%-30s %s\n", color.String(key), value)
	}
}
