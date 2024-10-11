package linux

import (
	"github.com/go-xuan/quanx/os/flagx"

	"quanx_tools/command"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Linux, "linux命令",
		flagx.StringOption("config", "配置文件", "gen.yaml"),
		flagx.BoolOption("check", "模板检测", false),
	).SetHandler(linux)
}

func linux() error {
	return nil
}
