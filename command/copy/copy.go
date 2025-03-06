package copy

import (
	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/base/ipx"

	"quanx_tools/command"
	"quanx_tools/common/utils"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Copy, "复制",
		flagx.BoolOption("bljmm", "堡垒机密码（蓝湖SSH登录）", false),
		flagx.BoolOption("ip", "本机IP", false),
	).SetExecutor(executor)
}

func executor() error {
	var text string
	if Command.GetOptionValue("bljmm").Bool() {
		text = "371ADDd70c27_"
	} else if Command.GetOptionValue("ip").Bool() {
		text = ipx.GetLocalIP()
	} else {
		Command.OptionsHelp()
	}
	if err := utils.WriteToClipboard(text); err != nil {
		return errorx.Wrap(err, "copy value to be pasted failed")
	}
	fmtx.Magenta.XPrintf("当前值%s已复制到粘贴板\n", text)
	return nil
}
