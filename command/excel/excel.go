package excel

import (
	"github.com/go-xuan/quanx/os/flagx"

	"quanx_tools/command"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Excel, "excel工具",
		flagx.StringOption("parse", "待解析excel文件", "parse.xlsx"),
		flagx.StringOption("mapping", "表头映射文件", "header_mapping.yaml"),
	).SetHandler(handler)
}

func handler() error {

	return nil
}
