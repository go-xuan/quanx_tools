package readme

import (
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/enumx"

	"quanx_tools/command"
	"quanx_tools/common/enums"
)

var (
	Command   = flagx.NewCommand(command.Readme, "README")
	LinuxEnum = enumx.NewStringEnum[string]()
	MacEnum   = enumx.NewStringEnum[string]()
)

func init() {
	Command.AddOption(
		flagx.BoolOption("linux", "linux常用命令", false),
		flagx.BoolOption("mac", "mac常用命令", false),
	).SetExecutor(executor)

	LinuxEnum.
		Add("", "")
	
	MacEnum.
		Add(`ifconfig en0 | grep "inet " | awk '{print $2}'`, "查看本机IP")
}

func executor() error {
	if Command.GetOptionValue("linux").Bool() {
		enums.Print(fmtx.Green, LinuxEnum)
	} else if Command.GetOptionValue("mac").Bool() {
		enums.Print(fmtx.Green, MacEnum)
	} else {
		Command.OptionsHelp()
	}
	return nil
}
