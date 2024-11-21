package readme

import (
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"

	"quanx_tools/command"
	"quanx_tools/common/enums"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Readme, "README",
		flagx.BoolOption("password", "账号密码", false),
		flagx.BoolOption("linux", "linux常用命令", false),
		flagx.BoolOption("mac", "mac常用命令", false),
	).SetExecutor(executor)
}

func executor() error {
	if Command.GetOptionValue("password").Bool() {
		enums.Print(fmtx.Green, enums.PasswordEnum)
	} else if Command.GetOptionValue("linux").Bool() {
		enums.Print(fmtx.Green, enums.LinuxEnum)
	} else if Command.GetOptionValue("mac").Bool() {
		enums.Print(fmtx.Green, enums.MacEnum)
	} else {
		Command.OptionsHelp()
	}
	return nil
}
