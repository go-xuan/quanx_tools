package main

import (
	"fmt"

	"github.com/go-xuan/quanx/os/flagx"

	"quanx_tools/command/copy"
	"quanx_tools/command/encrypt"
	"quanx_tools/command/gen"
	"quanx_tools/command/pwd"
	"quanx_tools/command/qr_code"
	"quanx_tools/command/rand"
	"quanx_tools/command/readme"
	"quanx_tools/command/redis"
	"quanx_tools/command/sql_fmt"
	"quanx_tools/command/time"
	"quanx_tools/command/tts"
)

func main() {
	flagx.Register(
		gen.Command,
		rand.Command,
		encrypt.Command,
		redis.Command,
		readme.Command,
		pwd.Command,
		//restic.Command,
		tts.Command,
		qr_code.Command,
		sql_fmt.Command,
		time.Command,
		copy.Command,
	)
	if err := flagx.Execute(); err != nil {
		fmt.Println(err)
	}
}
