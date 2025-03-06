package enums

import (
	"fmt"

	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/types/enumx"
)

func Print(color fmtx.Color, enum *enumx.StringEnum[string]) {
	for _, k := range enum.Keys() {
		fmt.Printf("%-30s %s\n", color.String(k), enum.Get(k))
	}
}
